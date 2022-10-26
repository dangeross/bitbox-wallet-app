// Copyright 2018 Shift Devices AG
// Copyright 2020 Shift Crypto AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/digitalbitbox/bitbox-wallet-app/backend"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/accounts"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/banners"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc"
	accountHandlers "github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/handlers"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/util"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/coin"
	coinpkg "github.com/digitalbitbox/bitbox-wallet-app/backend/coins/coin"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/eth"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/config"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/devices/bitbox"
	bitboxHandlers "github.com/digitalbitbox/bitbox-wallet-app/backend/devices/bitbox/handlers"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/devices/bitbox02"
	bitbox02Handlers "github.com/digitalbitbox/bitbox-wallet-app/backend/devices/bitbox02/handlers"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/devices/bitbox02bootloader"
	bitbox02bootloaderHandlers "github.com/digitalbitbox/bitbox-wallet-app/backend/devices/bitbox02bootloader/handlers"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/devices/device"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/exchanges"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/keystore"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/rates"
	utilConfig "github.com/digitalbitbox/bitbox-wallet-app/util/config"
	"github.com/digitalbitbox/bitbox-wallet-app/util/errp"
	"github.com/digitalbitbox/bitbox-wallet-app/util/jsonp"
	"github.com/digitalbitbox/bitbox-wallet-app/util/locker"
	"github.com/digitalbitbox/bitbox-wallet-app/util/logging"
	"github.com/digitalbitbox/bitbox-wallet-app/util/observable"
	"github.com/digitalbitbox/bitbox-wallet-app/util/socksproxy"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
)

// Backend models the API of the backend.
type Backend interface {
	observable.Interface

	Config() *config.Config
	DefaultAppConfig() config.AppConfig
	Coin(coinpkg.Code) (coinpkg.Coin, error)
	Testing() bool
	Accounts() []accounts.Interface
	Keystore() keystore.Keystore
	OnAccountInit(f func(accounts.Interface))
	OnAccountUninit(f func(accounts.Interface))
	OnDeviceInit(f func(device.Interface))
	OnDeviceUninit(f func(deviceID string))
	DevicesRegistered() map[string]device.Interface
	Start() <-chan interface{}
	DeregisterKeystore()
	Register(device device.Interface) error
	Deregister(deviceID string)
	RatesUpdater() *rates.RateUpdater
	DownloadCert(string) (string, error)
	CheckElectrumServer(*config.ServerInfo) error
	RegisterTestKeystore(string)
	NotifyUser(string)
	SystemOpen(string) error
	ReinitializeAccounts()
	CheckForUpdateIgnoringErrors() *backend.UpdateFile
	Banners() *banners.Banners
	Environment() backend.Environment
	ChartData() (*backend.Chart, error)
	SupportedCoins(keystore.Keystore) []coinpkg.Code
	CanAddAccount(coinpkg.Code, keystore.Keystore) (string, bool)
	CreateAndPersistAccountConfig(coinCode coinpkg.Code, name string, keystore keystore.Keystore) (accounts.Code, error)
	SetAccountActive(accountCode accounts.Code, active bool) error
	SetTokenActive(accountCode accounts.Code, tokenCode string, active bool) error
	RenameAccount(accountCode accounts.Code, name string) error
	AOPP() backend.AOPP
	AOPPCancel()
	AOPPApprove()
	AOPPChooseAccount(code accounts.Code)
}

// Handlers provides a web api to the backend.
type Handlers struct {
	Router  *mux.Router
	backend Backend
	// apiData consists of the port on which this API will run and the authorization token, generated by the
	// backend to secure the API call. The data is fed into the static javascript app
	// that is served, so the client knows where and how to connect to.
	apiData           *ConnectionData
	backendEvents     chan interface{}
	websocketUpgrader websocket.Upgrader
	log               *logrus.Entry
}

// ConnectionData contains the port and authorization token for communication with the backend.
type ConnectionData struct {
	port    int
	token   string
	devMode bool
}

// NewConnectionData creates a connection data struct which holds the port and token for the API.
// If the port is -1 or the token is empty, we assume dev-mode.
func NewConnectionData(port int, token string) *ConnectionData {
	return &ConnectionData{
		port:    port,
		token:   token,
		devMode: len(token) == 0,
	}
}

func (connectionData *ConnectionData) isDev() bool {
	return connectionData.port == -1 || connectionData.token == ""
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(
	backend Backend,
	connData *ConnectionData,
) *Handlers {
	log := logging.Get().WithGroup("handlers")
	router := mux.NewRouter()
	handlers := &Handlers{
		Router:        router,
		backend:       backend,
		apiData:       connData,
		backendEvents: make(chan interface{}, 1000),
		websocketUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		log: logging.Get().WithGroup("handlers"),
	}

	getAPIRouter := func(subrouter *mux.Router) func(string, func(*http.Request) (interface{}, error)) *mux.Route {
		return func(path string, f func(*http.Request) (interface{}, error)) *mux.Route {
			return subrouter.Handle(path, ensureAPITokenValid(handlers.apiMiddleware(connData.isDev(), f),
				connData, log))
		}
	}

	apiRouter := router.PathPrefix("/api").Subrouter()
	getAPIRouter(apiRouter)("/qr", handlers.getQRCodeHandler).Methods("GET")
	getAPIRouter(apiRouter)("/config", handlers.getAppConfigHandler).Methods("GET")
	getAPIRouter(apiRouter)("/config/default", handlers.getDefaultConfigHandler).Methods("GET")
	getAPIRouter(apiRouter)("/config", handlers.postAppConfigHandler).Methods("POST")
	getAPIRouter(apiRouter)("/native-locale", handlers.getNativeLocaleHandler).Methods("GET")
	getAPIRouter(apiRouter)("/notify-user", handlers.postNotifyHandler).Methods("POST")
	getAPIRouter(apiRouter)("/open", handlers.postOpenHandler).Methods("POST")
	getAPIRouter(apiRouter)("/update", handlers.getUpdateHandler).Methods("GET")
	getAPIRouter(apiRouter)("/banners/{key}", handlers.getBannersHandler).Methods("GET")
	getAPIRouter(apiRouter)("/using-mobile-data", handlers.getUsingMobileDataHandler).Methods("GET")
	getAPIRouter(apiRouter)("/version", handlers.getVersionHandler).Methods("GET")
	getAPIRouter(apiRouter)("/testing", handlers.getTestingHandler).Methods("GET")
	getAPIRouter(apiRouter)("/account-add", handlers.postAddAccountHandler).Methods("POST")
	getAPIRouter(apiRouter)("/keystores", handlers.getKeystoresHandler).Methods("GET")
	getAPIRouter(apiRouter)("/accounts", handlers.getAccountsHandler).Methods("GET")
	getAPIRouter(apiRouter)("/accounts/total-balance", handlers.getAccountsTotalBalanceHandler).Methods("GET")
	getAPIRouter(apiRouter)("/set-account-active", handlers.postSetAccountActiveHandler).Methods("POST")
	getAPIRouter(apiRouter)("/set-token-active", handlers.postSetTokenActiveHandler).Methods("POST")
	getAPIRouter(apiRouter)("/rename-account", handlers.postRenameAccountHandler).Methods("POST")
	getAPIRouter(apiRouter)("/accounts/reinitialize", handlers.postAccountsReinitializeHandler).Methods("POST")
	getAPIRouter(apiRouter)("/export-account-summary", handlers.postExportAccountSummary).Methods("POST")
	getAPIRouter(apiRouter)("/account-summary", handlers.getAccountSummary).Methods("GET")
	getAPIRouter(apiRouter)("/supported-coins", handlers.getSupportedCoinsHandler).Methods("GET")
	getAPIRouter(apiRouter)("/test/register", handlers.postRegisterTestKeystoreHandler).Methods("POST")
	getAPIRouter(apiRouter)("/test/deregister", handlers.postDeregisterTestKeystoreHandler).Methods("POST")
	getAPIRouter(apiRouter)("/rates", handlers.getRatesHandler).Methods("GET")
	getAPIRouter(apiRouter)("/coins/convertToPlainFiat", handlers.getConvertToPlainFiatHandler).Methods("GET")
	getAPIRouter(apiRouter)("/coins/convertFromFiat", handlers.getConvertFromFiatHandler).Methods("GET")
	getAPIRouter(apiRouter)("/coins/tltc/headers/status", handlers.getHeadersStatus(coinpkg.CodeTLTC)).Methods("GET")
	getAPIRouter(apiRouter)("/coins/tbtc/headers/status", handlers.getHeadersStatus(coinpkg.CodeTBTC)).Methods("GET")
	getAPIRouter(apiRouter)("/coins/ltc/headers/status", handlers.getHeadersStatus(coinpkg.CodeLTC)).Methods("GET")
	getAPIRouter(apiRouter)("/coins/btc/headers/status", handlers.getHeadersStatus(coinpkg.CodeBTC)).Methods("GET")
	getAPIRouter(apiRouter)("/coins/btc/set-unit", handlers.postBtcFormatUnit).Methods("POST")
	getAPIRouter(apiRouter)("/certs/download", handlers.postCertsDownloadHandler).Methods("POST")
	getAPIRouter(apiRouter)("/electrum/check", handlers.postElectrumCheckHandler).Methods("POST")
	getAPIRouter(apiRouter)("/socksproxy/check", handlers.postSocksProxyCheck).Methods("POST")
	getAPIRouter(apiRouter)("/exchange/moonpay/buy-supported/{code}", handlers.getExchangeMoonpayBuySupported).Methods("GET")
	getAPIRouter(apiRouter)("/exchange/moonpay/buy/{code}", handlers.getExchangeMoonpayBuy).Methods("GET")
	getAPIRouter(apiRouter)("/aopp", handlers.getAOPPHandler).Methods("GET")
	getAPIRouter(apiRouter)("/aopp/cancel", handlers.postAOPPCancelHandler).Methods("POST")
	getAPIRouter(apiRouter)("/aopp/approve", handlers.postAOPPApproveHandler).Methods("POST")
	getAPIRouter(apiRouter)("/aopp/choose-account", handlers.postAOPPChooseAccountHandler).Methods("POST")

	devicesRouter := getAPIRouter(apiRouter.PathPrefix("/devices").Subrouter())
	devicesRouter("/registered", handlers.getDevicesRegisteredHandler).Methods("GET")

	handlersMapLock := locker.Locker{}

	accountHandlersMap := map[accounts.Code]*accountHandlers.Handlers{}
	getAccountHandlers := func(accountCode accounts.Code) *accountHandlers.Handlers {
		defer handlersMapLock.Lock()()
		if _, ok := accountHandlersMap[accountCode]; !ok {
			accountHandlersMap[accountCode] = accountHandlers.NewHandlers(getAPIRouter(
				apiRouter.PathPrefix(fmt.Sprintf("/account/%s", accountCode)).Subrouter(),
			), log)
		}
		accHandlers := accountHandlersMap[accountCode]
		log.WithField("account-handlers", accHandlers).Debug("Account handlers")
		return accHandlers
	}

	backend.OnAccountInit(func(account accounts.Interface) {
		log.WithField("code", account.Config().Code).Debug("Initializing account")
		getAccountHandlers(account.Config().Code).Init(account)
	})
	backend.OnAccountUninit(func(account accounts.Interface) {
		getAccountHandlers(account.Config().Code).Uninit()
	})

	deviceHandlersMap := map[string]*bitboxHandlers.Handlers{}
	getDeviceHandlers := func(deviceID string) *bitboxHandlers.Handlers {
		defer handlersMapLock.Lock()()
		if _, ok := deviceHandlersMap[deviceID]; !ok {
			deviceHandlersMap[deviceID] = bitboxHandlers.NewHandlers(getAPIRouter(
				apiRouter.PathPrefix(fmt.Sprintf("/devices/%s", deviceID)).Subrouter(),
			), log)
		}
		return deviceHandlersMap[deviceID]
	}

	bitbox02HandlersMap := map[string]*bitbox02Handlers.Handlers{}
	getBitBox02Handlers := func(deviceID string) *bitbox02Handlers.Handlers {
		defer handlersMapLock.Lock()()
		if _, ok := bitbox02HandlersMap[deviceID]; !ok {
			bitbox02HandlersMap[deviceID] = bitbox02Handlers.NewHandlers(getAPIRouter(
				apiRouter.PathPrefix(fmt.Sprintf("/devices/bitbox02/%s", deviceID)).Subrouter(),
			), log)
		}
		return bitbox02HandlersMap[deviceID]
	}

	bitbox02BootloaderHandlersMap := map[string]*bitbox02bootloaderHandlers.Handlers{}
	getBitBox02BootloaderHandlers := func(deviceID string) *bitbox02bootloaderHandlers.Handlers {
		defer handlersMapLock.Lock()()
		if _, ok := bitbox02BootloaderHandlersMap[deviceID]; !ok {
			bitbox02BootloaderHandlersMap[deviceID] = bitbox02bootloaderHandlers.NewHandlers(getAPIRouter(
				apiRouter.PathPrefix(fmt.Sprintf("/devices/bitbox02-bootloader/%s", deviceID)).Subrouter(),
			), log)
		}
		return bitbox02BootloaderHandlersMap[deviceID]
	}

	backend.OnDeviceInit(func(device device.Interface) {
		switch specificDevice := device.(type) {
		case *bitbox.Device:
			getDeviceHandlers(device.Identifier()).Init(specificDevice)
		case *bitbox02.Device:
			getBitBox02Handlers(device.Identifier()).Init(specificDevice)
		case *bitbox02bootloader.Device:
			getBitBox02BootloaderHandlers(device.Identifier()).Init(specificDevice)
		}
	})
	backend.OnDeviceUninit(func(deviceID string) {
		getDeviceHandlers(deviceID).Uninit()
	})

	apiRouter.HandleFunc("/events", handlers.eventsHandler)

	// The backend relays events in two ways:
	// a) old school through the channel returned by Start()
	// b) new school via observable.
	// Merge both.
	events := backend.Start()
	go func() {
		for {
			handlers.backendEvents <- <-events
		}
	}()
	backend.Observe(func(event observable.Event) { handlers.backendEvents <- event })

	return handlers
}

// Events returns the push notifications channel.
func (handlers *Handlers) Events() <-chan interface{} {
	return handlers.backendEvents
}

func writeJSON(w io.Writer, value interface{}) {
	if err := json.NewEncoder(w).Encode(value); err != nil {
		panic(err)
	}
}

type activeToken struct {
	// TokenCode is the token code as defined in erc20.go, e.g. "eth-erc20-usdt".
	TokenCode string `json:"tokenCode"`
	// AccountCode is the code of the account, which is not the same as the TokenCode, as there can
	// be many accounts for the same token.
	AccountCode accounts.Code `json:"accountCode"`
}

type accountJSON struct {
	Active                bool          `json:"active"`
	CoinCode              coinpkg.Code  `json:"coinCode"`
	CoinUnit              string        `json:"coinUnit"`
	CoinName              string        `json:"coinName"`
	Code                  accounts.Code `json:"code"`
	Name                  string        `json:"name"`
	IsToken               bool          `json:"isToken"`
	ActiveTokens          []activeToken `json:"activeTokens,omitempty"`
	BlockExplorerTxPrefix string        `json:"blockExplorerTxPrefix"`
}

func newAccountJSON(account accounts.Interface, activeTokens []activeToken) *accountJSON {
	eth, ok := account.Coin().(*eth.Coin)
	isToken := ok && eth.ERC20Token() != nil
	return &accountJSON{
		Active:                account.Config().Active,
		CoinCode:              account.Coin().Code(),
		CoinUnit:              account.Coin().Unit(false),
		CoinName:              account.Coin().Name(),
		Code:                  account.Config().Code,
		Name:                  account.Config().Name,
		IsToken:               isToken,
		ActiveTokens:          activeTokens,
		BlockExplorerTxPrefix: account.Coin().BlockExplorerTransactionURLPrefix(),
	}
}

func (handlers *Handlers) getQRCodeHandler(r *http.Request) (interface{}, error) {
	data := r.URL.Query().Get("data")
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, errp.WithStack(err)
	}
	bytes, err := qr.PNG(256)
	if err != nil {
		return nil, errp.WithStack(err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes), nil
}

func (handlers *Handlers) getAppConfigHandler(_ *http.Request) (interface{}, error) {
	return handlers.backend.Config().AppConfig(), nil
}

func (handlers *Handlers) getDefaultConfigHandler(_ *http.Request) (interface{}, error) {
	return handlers.backend.DefaultAppConfig(), nil
}

func (handlers *Handlers) postAppConfigHandler(r *http.Request) (interface{}, error) {
	appConfig := config.AppConfig{}
	if err := json.NewDecoder(r.Body).Decode(&appConfig); err != nil {
		return nil, errp.WithStack(err)
	}
	return nil, handlers.backend.Config().SetAppConfig(appConfig)
}

// getNativeLocaleHandler returns user preferred UI language as reported
// by the native app layer.
// The response value may be invalid or unsupported by the app.
func (handlers *Handlers) getNativeLocaleHandler(*http.Request) (interface{}, error) {
	return handlers.backend.Environment().NativeLocale(), nil
}

func (handlers *Handlers) postNotifyHandler(r *http.Request) (interface{}, error) {
	payload := struct {
		Text string `json:"text"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, errp.WithStack(err)
	}
	handlers.backend.NotifyUser(payload.Text)
	return nil, nil
}

func (handlers *Handlers) postOpenHandler(r *http.Request) (interface{}, error) {
	var url string
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return nil, errp.WithStack(err)
	}
	return nil, handlers.backend.SystemOpen(url)
}

func (handlers *Handlers) getUpdateHandler(_ *http.Request) (interface{}, error) {
	return handlers.backend.CheckForUpdateIgnoringErrors(), nil
}

func (handlers *Handlers) getBannersHandler(r *http.Request) (interface{}, error) {
	return handlers.backend.Banners().GetMessage(banners.MessageKey(mux.Vars(r)["key"])), nil
}

func (handlers *Handlers) getUsingMobileDataHandler(r *http.Request) (interface{}, error) {
	return handlers.backend.Environment().UsingMobileData(), nil
}

func (handlers *Handlers) getVersionHandler(_ *http.Request) (interface{}, error) {
	return backend.Version.String(), nil
}

func (handlers *Handlers) getTestingHandler(_ *http.Request) (interface{}, error) {
	return handlers.backend.Testing(), nil
}

func (handlers *Handlers) postAddAccountHandler(r *http.Request) (interface{}, error) {
	var jsonBody struct {
		CoinCode coinpkg.Code `json:"coinCode"`
		Name     string       `json:"name"`
	}

	type response struct {
		Success      bool          `json:"success"`
		AccountCode  accounts.Code `json:"accountCode,omitempty"`
		ErrorMessage string        `json:"errorMessage,omitempty"`
		ErrorCode    string        `json:"errorCode,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil {
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}

	keystore := handlers.backend.Keystore()
	if keystore == nil {
		return response{Success: false, ErrorMessage: "Keystore not found"}, nil
	}

	accountCode, err := handlers.backend.CreateAndPersistAccountConfig(jsonBody.CoinCode, jsonBody.Name, keystore)
	if err != nil {
		handlers.log.WithError(err).Error("Could not add account")
		if errCode, ok := errp.Cause(err).(backend.ErrorCode); ok {
			return response{Success: false, ErrorCode: string(errCode)}, nil
		}
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}
	return response{Success: true, AccountCode: accountCode}, nil
}

func (handlers *Handlers) getKeystoresHandler(_ *http.Request) (interface{}, error) {
	type json struct {
		Type keystore.Type `json:"type"`
	}
	keystores := []*json{}

	keystore := handlers.backend.Keystore()
	if keystore != nil {
		keystores = append(keystores, &json{
			Type: keystore.Type(),
		})
	}
	return keystores, nil
}

func (handlers *Handlers) getAccountsHandler(_ *http.Request) (interface{}, error) {
	accounts := []*accountJSON{}
	persistedAccounts := handlers.backend.Config().AccountsConfig()
	for _, account := range handlers.backend.Accounts() {
		var activeTokens []activeToken
		if account.Coin().Code() == coinpkg.CodeETH {
			persistedAccount := persistedAccounts.Lookup(account.Config().Code)
			if persistedAccount == nil {
				handlers.log.WithField("code", account.Config().Code).Error("account not found in accounts database")
				continue
			}
			for _, tokenCode := range persistedAccount.ActiveTokens {
				activeTokens = append(activeTokens, activeToken{
					TokenCode:   tokenCode,
					AccountCode: backend.Erc20AccountCode(account.Config().Code, tokenCode),
				})
			}
		}
		accounts = append(accounts, newAccountJSON(account, activeTokens))
	}
	return accounts, nil
}

func (handlers *Handlers) postBtcFormatUnit(r *http.Request) (interface{}, error) {
	type response struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"errorMessage,omitempty"`
	}

	var request struct {
		Unit coinpkg.BtcUnit `json:"unit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return response{Success: false}, nil
	}

	unit := request.Unit

	// update BTC format unit for Coins
	btcCoin, err := handlers.backend.Coin(coinpkg.CodeBTC)
	if err != nil {
		return response{Success: false}, nil
	}
	btcCoin.(*btc.Coin).SetFormatUnit(unit)

	btcCoin, err = handlers.backend.Coin(coinpkg.CodeTBTC)
	if err != nil {
		return response{Success: false}, nil
	}
	btcCoin.(*btc.Coin).SetFormatUnit(unit)

	// update BTC format unit for fiat conversions
	for _, account := range handlers.backend.Accounts() {
		account.Config().BtcCurrencyUnit = unit
	}

	return response{Success: true}, nil
}

func (handlers *Handlers) getAccountsTotalBalanceHandler(_ *http.Request) (interface{}, error) {
	totalPerCoin := make(map[coin.Code]*big.Int)
	conversionsPerCoin := make(map[coin.Code]map[string]string)

	totalAmount := make(map[coin.Code]accountHandlers.FormattedAmount)

	for _, account := range handlers.backend.Accounts() {
		if !account.Config().Active {
			continue
		}
		if account.FatalError() {
			continue
		}
		err := account.Initialize()
		if err != nil {
			return nil, err
		}
		coinCode := account.Coin().Code()
		b, err := account.Balance()
		if err != nil {
			return nil, err
		}
		amount := b.Available()
		if _, ok := totalPerCoin[coinCode]; !ok {
			totalPerCoin[coinCode] = amount.BigInt()

		} else {
			totalPerCoin[coinCode] = new(big.Int).Add(totalPerCoin[coinCode], amount.BigInt())
		}

		conversionsPerCoin[coinCode] = coin.Conversions(
			coin.NewAmount(totalPerCoin[coinCode]),
			account.Coin(),
			false,
			account.Config().RateUpdater,
			util.FormatBtcAsSat(handlers.backend.Config().AppConfig().Backend.BtcUnit))
	}

	for k, v := range totalPerCoin {
		currentCoin, err := handlers.backend.Coin(k)
		if err != nil {
			return nil, err
		}
		totalAmount[k] = accountHandlers.FormattedAmount{
			Amount:      currentCoin.FormatAmount(coin.NewAmount(v), false),
			Unit:        currentCoin.GetFormatUnit(),
			Conversions: conversionsPerCoin[k],
		}
	}
	return totalAmount, nil
}

func (handlers *Handlers) postSetAccountActiveHandler(r *http.Request) (interface{}, error) {
	var jsonBody struct {
		AccountCode accounts.Code `json:"accountCode"`
		Active      bool          `json:"active"`
	}

	type response struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"errorMessage,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil {
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}
	if err := handlers.backend.SetAccountActive(jsonBody.AccountCode, jsonBody.Active); err != nil {
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}
	return response{Success: true}, nil
}

func (handlers *Handlers) postSetTokenActiveHandler(r *http.Request) (interface{}, error) {
	var jsonBody struct {
		AccountCode accounts.Code `json:"accountCode"`
		TokenCode   string        `json:"tokenCode"`
		Active      bool          `json:"active"`
	}

	type response struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"errorMessage,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil {
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}
	if err := handlers.backend.SetTokenActive(jsonBody.AccountCode, jsonBody.TokenCode, jsonBody.Active); err != nil {
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}
	return response{Success: true}, nil
}

func (handlers *Handlers) postRenameAccountHandler(r *http.Request) (interface{}, error) {
	var jsonBody struct {
		AccountCode accounts.Code `json:"accountCode"`
		Name        string        `json:"name"`
	}

	type response struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"errorMessage,omitempty"`
		ErrorCode    string `json:"errorCode,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil {
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}
	if err := handlers.backend.RenameAccount(jsonBody.AccountCode, jsonBody.Name); err != nil {
		if errCode, ok := errp.Cause(err).(backend.ErrorCode); ok {
			return response{Success: false, ErrorCode: string(errCode)}, nil
		}
		return response{Success: false, ErrorMessage: err.Error()}, nil
	}
	return response{Success: true}, nil
}

func (handlers *Handlers) postAccountsReinitializeHandler(_ *http.Request) (interface{}, error) {
	handlers.backend.ReinitializeAccounts()
	return nil, nil
}

func (handlers *Handlers) getDevicesRegisteredHandler(_ *http.Request) (interface{}, error) {
	jsonDevices := map[string]string{}
	for deviceID, device := range handlers.backend.DevicesRegistered() {
		jsonDevices[deviceID] = device.ProductName()
	}
	return jsonDevices, nil
}

func (handlers *Handlers) postRegisterTestKeystoreHandler(r *http.Request) (interface{}, error) {
	if !handlers.backend.Testing() {
		return nil, errp.New("Test keystore not available")
	}
	jsonBody := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil {
		return nil, errp.WithStack(err)
	}
	pin := jsonBody["pin"]
	handlers.backend.RegisterTestKeystore(pin)
	return nil, nil
}

func (handlers *Handlers) postDeregisterTestKeystoreHandler(_ *http.Request) (interface{}, error) {
	handlers.backend.DeregisterKeystore()
	return nil, nil
}

func (handlers *Handlers) getRatesHandler(_ *http.Request) (interface{}, error) {
	return handlers.backend.RatesUpdater().LatestPrice(), nil
}

func (handlers *Handlers) getConvertToPlainFiatHandler(r *http.Request) (interface{}, error) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amount := r.URL.Query().Get("amount")

	currentCoin, err := handlers.backend.Coin(coinpkg.Code(from))
	if err != nil {
		logrus.Error(err.Error())
		return map[string]interface{}{
			"success": false,
		}, nil
	}

	coinAmount, err := currentCoin.ParseAmount(amount)
	if err != nil {
		logrus.Error(err.Error())
		return map[string]interface{}{
			"success": false,
		}, nil
	}

	coinUnitAmount := new(big.Rat).SetFloat64(currentCoin.ToUnit(coinAmount, false))

	unit := currentCoin.Unit(false)
	rate := handlers.backend.RatesUpdater().LatestPrice()[unit][to]

	convertedAmount := new(big.Rat).Mul(coinUnitAmount, new(big.Rat).SetFloat64(rate))

	btcUnit := handlers.backend.Config().AppConfig().Backend.BtcUnit
	return map[string]interface{}{
		"success":    true,
		"fiatAmount": coinpkg.FormatAsPlainCurrency(convertedAmount, to, util.FormatBtcAsSat(btcUnit)),
	}, nil
}

func (handlers *Handlers) getConvertFromFiatHandler(r *http.Request) (interface{}, error) {
	isFee := false
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	currentCoin, err := handlers.backend.Coin(coinpkg.Code(to))
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"errMsg":  "internal error",
		}, nil
	}

	fiatStr := r.URL.Query().Get("amount")
	fiatRat, valid := new(big.Rat).SetString(fiatStr)
	if !valid {
		return map[string]interface{}{
			"success": false,
			"errMsg":  "invalid amount",
		}, nil
	}

	unit := currentCoin.Unit(isFee)
	switch unit { // HACK: fake rates for testnet coins
	case "TBTC", "TLTC", "TETH", "RETH":
		unit = unit[1:]
	case "GOETH":
		unit = unit[2:]
	}

	if from == rates.BTC.String() && handlers.backend.Config().AppConfig().Backend.BtcUnit == coinpkg.BtcUnitSats {
		fiatRat = coinpkg.Sat2Btc(fiatRat)
	}

	rate := handlers.backend.RatesUpdater().LatestPrice()[unit][from]
	result := coin.NewAmountFromInt64(0)
	if rate != 0.0 {
		amountRat := new(big.Rat).Quo(fiatRat, new(big.Rat).SetFloat64(rate))
		result = currentCoin.SetAmount(amountRat, false)
	}
	return map[string]interface{}{
		"success": true,
		"amount":  currentCoin.FormatAmount(result, false),
	}, nil
}

func (handlers *Handlers) getHeadersStatus(coinCode coinpkg.Code) func(*http.Request) (interface{}, error) {
	return func(_ *http.Request) (interface{}, error) {
		coin, err := handlers.backend.Coin(coinCode)
		if err != nil {
			return nil, err
		}
		return coin.(*btc.Coin).Headers().Status()
	}
}

func (handlers *Handlers) postCertsDownloadHandler(r *http.Request) (interface{}, error) {
	var server string
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		return nil, errp.WithStack(err)
	}
	pemCert, err := handlers.backend.DownloadCert(server)
	if err != nil {
		return map[string]interface{}{
			"success":      false,
			"errorMessage": err.Error(),
		}, nil
	}
	return map[string]interface{}{
		"success": true,
		"pemCert": pemCert,
	}, nil
}

func (handlers *Handlers) postElectrumCheckHandler(r *http.Request) (interface{}, error) {
	var serverInfo config.ServerInfo
	if err := json.NewDecoder(r.Body).Decode(&serverInfo); err != nil {
		return nil, errp.WithStack(err)
	}

	if err := handlers.backend.CheckElectrumServer(&serverInfo); err != nil {
		return map[string]interface{}{
			"success":      false,
			"errorMessage": err.Error(),
		}, nil
	}
	return map[string]interface{}{
		"success": true,
	}, nil
}

func (handlers *Handlers) postSocksProxyCheck(r *http.Request) (interface{}, error) {
	var endpoint string
	if err := json.NewDecoder(r.Body).Decode(&endpoint); err != nil {
		return nil, errp.WithStack(err)
	}

	type response struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"errorMessage,omitempty"`
	}

	err := socksproxy.NewSocksProxy(true, endpoint).Validate()
	if err != nil {
		return response{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	return response{
		Success: true,
	}, nil
}

func (handlers *Handlers) eventsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := handlers.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	sendChan, quitChan := runWebsocket(conn, handlers.apiData, handlers.log)
	go func() {
		for {
			select {
			case <-quitChan:
				return
			default:
				select {
				case <-quitChan:
					return
				case event := <-handlers.backendEvents:
					sendChan <- jsonp.MustMarshal(event)
				}
			}
		}
	}()
}

// isAPITokenValid checks whether we are in dev or prod mode and, if we are in prod mode, verifies
// that an authorization token is received as an HTTP Authorization header and that it is valid.
func isAPITokenValid(w http.ResponseWriter, r *http.Request, apiData *ConnectionData, log *logrus.Entry) bool {
	methodLogEntry := log.
		WithField("path", r.URL.Path).
		WithField("method", r.Method)
	methodLogEntry.Debug("endpoint")
	// In dev mode, we allow unauthorized requests
	if apiData.devMode {
		return true
	}

	if len(r.Header.Get("Authorization")) == 0 {
		methodLogEntry.Error("Missing token in API request. WARNING: this could be an attack on the API")
		http.Error(w, "missing token "+r.URL.Path, http.StatusUnauthorized)
		return false
	} else if len(r.Header.Get("Authorization")) != 0 && r.Header.Get("Authorization") != "Basic "+apiData.token {
		methodLogEntry.Error("Incorrect token in API request. WARNING: this could be an attack on the API")
		http.Error(w, "incorrect token", http.StatusUnauthorized)
		return false
	}
	return true
}

// ensureAPITokenValid wraps the given handler with another handler function that calls isAPITokenValid().
func ensureAPITokenValid(h http.Handler, apiData *ConnectionData, log *logrus.Entry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAPITokenValid(w, r, apiData, log) {
			h.ServeHTTP(w, r)
		}
	})
}

func (handlers *Handlers) apiMiddleware(devMode bool, h func(*http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// recover from all panics and log error before panicking again
			if r := recover(); r != nil {
				handlers.log.WithField("panic", true).Errorf("%v\n%s", r, string(debug.Stack()))
				writeJSON(w, map[string]string{"error": fmt.Sprintf("%v", r)})
			}
		}()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if devMode {
			// This enables us to run a server on a different port serving just the UI, while still
			// allowing it to access the API.
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		}
		value, err := h(r)
		if err != nil {
			handlers.log.WithError(err).Error("endpoint failed")
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, value)
	})
}
func (handlers *Handlers) getAccountSummary(_ *http.Request) (interface{}, error) {
	return handlers.backend.ChartData()
}

// getSupportedCoinsHandler returns an array of coin codes for which you can add an account.
// Exactly one keystore must be connected, otherwise an empty array is returned.
func (handlers *Handlers) getSupportedCoinsHandler(_ *http.Request) (interface{}, error) {
	type element struct {
		CoinCode             coinpkg.Code `json:"coinCode"`
		Name                 string       `json:"name"`
		CanAddAccount        bool         `json:"canAddAccount"`
		SuggestedAccountName string       `json:"suggestedAccountName"`
	}
	keystore := handlers.backend.Keystore()
	if keystore == nil {
		return []string{}, nil
	}
	var result []element
	for _, coinCode := range handlers.backend.SupportedCoins(keystore) {
		coin, err := handlers.backend.Coin(coinCode)
		if err != nil {
			continue
		}
		suggestedAccountName, canAddAccount := handlers.backend.CanAddAccount(coinCode, keystore)
		result = append(result, element{
			CoinCode:             coinCode,
			Name:                 coin.Name(),
			CanAddAccount:        canAddAccount,
			SuggestedAccountName: suggestedAccountName,
		})
	}
	return result, nil
}

func (handlers *Handlers) postExportAccountSummary(_ *http.Request) (interface{}, error) {
	name := time.Now().Format("2006-01-02-at-15-04-05-") + "Accounts-Summary.csv"
	downloadsDir, err := utilConfig.DownloadsDir()
	if err != nil {
		return nil, err
	}
	suggestedPath := filepath.Join(downloadsDir, name)
	path := handlers.backend.Environment().GetSaveFilename(suggestedPath)
	if path == "" {
		return nil, nil
	}
	handlers.log.Infof("Export account summary %s.", path)

	file, err := os.Create(path)
	if err != nil {
		return nil, errp.WithStack(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			handlers.log.WithError(err).Error("Could not close the account summary file.")
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{
		"Coin",
		"Name",
		"Balance",
		"Unit",
		"Type",
		"Xpubs",
	})
	if err != nil {
		return nil, errp.WithStack(err)
	}

	for _, account := range handlers.backend.Accounts() {
		if !account.Config().Active {
			continue
		}
		if account.FatalError() {
			continue
		}
		err := account.Initialize()
		if err != nil {
			return nil, err
		}
		coin := account.Coin().Code()
		accountName := account.Config().Name
		balance, err := account.Balance()
		if err != nil {
			return nil, err
		}
		unit := account.Coin().SmallestUnit()
		var accountType string
		var xpubs []string
		signingConfigurations := account.Info().SigningConfigurations
		accountType = "xpubs"
		for _, signingConfiguration := range signingConfigurations {
			xpubs = append(xpubs, signingConfiguration.ExtendedPublicKey().String())
		}

		err = writer.Write([]string{
			string(coin),
			accountName,
			balance.Available().BigInt().String(),
			unit,
			accountType,
			strings.Join(xpubs, "; "),
		})
		if err != nil {
			return nil, errp.WithStack(err)
		}
	}
	return path, nil
}

func (handlers *Handlers) getExchangeMoonpayBuySupported(r *http.Request) (interface{}, error) {
	acctCode := mux.Vars(r)["code"]
	// TODO: Refactor to make use of a map.
	var acct accounts.Interface
	for _, a := range handlers.backend.Accounts() {
		if !a.Config().Active {
			continue
		}
		if string(a.Config().Code) == acctCode {
			acct = a
			break
		}
	}
	// TODO: Offline() can be removed from here once there is a unified way of initializing accounts
	// and showing sync status, offline/fatal states, etc.
	return acct != nil && acct.Offline() == nil && exchanges.IsMoonpaySupported(acct.Coin().Code()), nil
}

func (handlers *Handlers) getAOPPHandler(r *http.Request) (interface{}, error) {
	return handlers.backend.AOPP(), nil
}

func (handlers *Handlers) getExchangeMoonpayBuy(r *http.Request) (interface{}, error) {
	acctCode := accounts.Code(mux.Vars(r)["code"])
	// TODO: Refactor to make use of a map.
	var acct accounts.Interface
	for _, a := range handlers.backend.Accounts() {
		if !a.Config().Active {
			continue
		}
		if a.Config().Code == acctCode {
			acct = a
			break
		}
	}
	if acct == nil {
		return nil, fmt.Errorf("unknown account code %q", acctCode)
	}

	if err := acct.Initialize(); err != nil {
		return nil, err
	}

	params := exchanges.BuyMoonpayParams{
		Fiat: handlers.backend.Config().AppConfig().Backend.MainFiat,
		Lang: handlers.backend.Config().AppConfig().Backend.UserLanguage,
	}
	buy, err := exchanges.BuyMoonpay(acct, params)
	if err != nil {
		return nil, err
	}
	resp := struct {
		URL     string `json:"url"`
		Address string `json:"address"`
	}{
		URL:     buy.URL,
		Address: buy.Address,
	}
	return resp, nil
}

func (handlers *Handlers) postAOPPChooseAccountHandler(r *http.Request) (interface{}, error) {
	var request struct {
		AccountCode accounts.Code `json:"accountCode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errp.WithStack(err)
	}

	handlers.backend.AOPPChooseAccount(request.AccountCode)
	return nil, nil
}

func (handlers *Handlers) postAOPPCancelHandler(r *http.Request) (interface{}, error) {
	handlers.backend.AOPPCancel()
	return nil, nil
}

func (handlers *Handlers) postAOPPApproveHandler(r *http.Request) (interface{}, error) {
	handlers.backend.AOPPApprove()
	return nil, nil
}
