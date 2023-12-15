/**
 * Copyright 2018 Shift Devices AG
 * Copyright 2023 Shift Crypto AG
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Column, ColumnButtons, Grid, GuideWrapper, GuidedContent, Header, Main } from '../../../components/layout';
import { useTranslation } from 'react-i18next';
import { View, ViewButtons, ViewContent } from '../../../components/view/view';
import { Button, Input, OptionalLabel } from '../../../components/forms';
import { ChangeEvent, useCallback, useEffect, useState } from 'react';
import {
  InputType,
  InputTypeVariant,
  LnUrlWithdrawResultVariant,
  OpenChannelFeeResponse,
  Payment,
  PaymentStatus,
  PaymentTypeFilter,
  ReceivePaymentResponse,
  SdkError,
  getListPayments,
  getOpenChannelFee,
  getParseInput,
  postLnurlWithdraw,
  postReceivePayment,
  subscribeListPayments
} from '../../../api/lightning';
import { route } from '../../../utils/route';
import { toMsat, toSat } from '../../../utils/conversion';
import { Status } from '../../../components/status/status';
import { QRCode } from '../../../components/qrcode/qrcode';
import { unsubscribe } from '../../../utils/subscriptions';
import { Spinner } from '../../../components/spinner/Spinner';
import { Checked, Copy, EditActive } from '../../../components/icon';
import styles from './receive.module.css';
import { ScanQRVideo } from '../../account/send/components/inputs/scan-qr-video';

type TStep = 'select-amount' | 'scan' | 'wait' | 'invoice' | 'success';

export function Receive() {
  const { t } = useTranslation();
  const [amountSats, setAmountSats] = useState<number>(0);
  const [amountSatsText, setAmountSatsText] = useState<string>('');
  const [description, setDescription] = useState<string>('');
  const [minWithdrawable, setMinWithdrawable] = useState<number>(0);
  const [maxWithdrawable, setMaxWithdrawable] = useState<number>(0);
  const [openChannelFeeResponse, setOpenChannelFeeResponse] = useState<OpenChannelFeeResponse>();
  const [parsedInput, setParsedInput] = useState<InputType>();
  const [rawInputError, setRawInputError] = useState<string>();
  const [receiveDisabled, setReceiveDisabled] = useState<boolean>();
  const [receivePaymentResponse, setReceivePaymentResponse] = useState<ReceivePaymentResponse>();
  const [receiveError, setReceiveError] = useState<string>();
  const [showOpenChannelWarning, setShowOpenChannelWarning] = useState<boolean>(false);
  const [step, setStep] = useState<TStep>('select-amount');
  const [payments, setPayments] = useState<Payment[]>();

  const back = () => {
    setReceiveError(undefined);
    switch (step) {
    case 'select-amount':
      route('/lightning');
      break;
    case 'scan':
    case 'invoice':
    case 'success':
      setStep('select-amount');
      if (step === 'success') {
        setAmountSatsText('');
      }
      break;
    }
  };

  const onAmountSatsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target as HTMLInputElement;
    setAmountSatsText(target.value);
  };

  const onDescriptionChange = (event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target as HTMLInputElement;
    setDescription(target.value);
  };

  const onPaymentsChange = useCallback(() => {
    getListPayments({ filters: [PaymentTypeFilter.RECEIVED], limit: 5 }).then((payments) => setPayments(payments));
  }, []);

  const onScanQrCode = (() => {
    setStep('scan');
  });

  useEffect(() => {
    const subscriptions = [subscribeListPayments(onPaymentsChange)];
    return () => unsubscribe(subscriptions);
  }, [onPaymentsChange]);

  useEffect(() => {
    setAmountSats(+amountSatsText);
  }, [amountSatsText]);

  useEffect(() => {
    (async () => {
      if (amountSats > 0) {
        const openChannelFeeResponse = await getOpenChannelFee({ amountMsat: toMsat(amountSats) });
        setOpenChannelFeeResponse(openChannelFeeResponse);
        setShowOpenChannelWarning(openChannelFeeResponse.feeMsat > 0);
      }
    })();
  }, [amountSats]);

  useEffect(() => {
    if (payments && receivePaymentResponse && step === 'invoice') {
      const payment = payments.find((payment) => payment.id === receivePaymentResponse.lnInvoice.paymentHash);
      if (payment?.status === PaymentStatus.COMPLETE) {
        setStep('success');
        setTimeout(() => route('/lightning'), 5000);
      }
    }
  }, [payments, receivePaymentResponse, step]);

  useEffect(() => {
    if (parsedInput && parsedInput.type === InputTypeVariant.LN_URL_WITHDRAW) {
      setReceiveDisabled(amountSats <= 0 || amountSats < minWithdrawable || amountSats > maxWithdrawable);
    } else {
      setReceiveDisabled(amountSats <= 0);
    }
  }, [amountSats, maxWithdrawable, minWithdrawable, parsedInput]);

  const parseInput = useCallback(async (rawInput: string) => {
    setRawInputError(undefined);
    try {
      const result = await getParseInput({ s: rawInput });
      switch (result.type) {
      case InputTypeVariant.LN_URL_WITHDRAW:
        setAmountSatsText('');
        setMaxWithdrawable(toSat(result.data.maxWithdrawable));
        setMinWithdrawable(toSat(result.data.minWithdrawable));
        setDescription(result.data.defaultDescription);
        setParsedInput(result);
        setStep('select-amount');
        break;
      default:
        setRawInputError('Invalid input');
      }
    } catch (e) {
      if (e instanceof SdkError) {
        setRawInputError(e.message);
      } else {
        setRawInputError(String(e));
      }
    }
  }, []);

  const receivePayment = async () => {
    setStep('wait');
    setReceiveError(undefined);
    try {
      if (parsedInput && parsedInput.type === InputTypeVariant.LN_URL_WITHDRAW) {
        const lnUrlWithdrawResult = await postLnurlWithdraw({
          data: parsedInput.data,
          amountMsat: toMsat(amountSats),
          description
        });
        switch (lnUrlWithdrawResult.type) {
        case LnUrlWithdrawResultVariant.ERROR_STATUS:
          setReceiveError(lnUrlWithdrawResult.data.reason);
          setStep('select-amount');
          break;
        case LnUrlWithdrawResultVariant.OK:
          setStep('success');
          setTimeout(() => route('/lightning'), 5000);
          break;
        }
      } else {
        const receivePaymentResponse = await postReceivePayment({
          amountMsat: toMsat(amountSats),
          description,
          openingFeeParams: openChannelFeeResponse?.usedFeeParams
        });
        setReceivePaymentResponse(receivePaymentResponse);
        setStep('invoice');
      }
    } catch (e) {
      if (e instanceof SdkError) {
        setReceiveError(e.message);
      } else {
        setReceiveError(String(e));
      }
    }
  };

  const renderSteps = () => {
    switch (step) {
    case 'select-amount':
      return (
        <View>
          <ViewContent>
            <Grid col="1">
              <Column>
                {!parsedInput && (<h1 className={styles.title}>{t('lightning.receive.subtitle')}</h1>)}
                <Input
                  type="number"
                  min="0"
                  label={t('lightning.receive.amountSats.label')}
                  labelSection={
                    parsedInput && (
                      <OptionalLabel>{t('lightning.receive.amountSats.limitLabel', { maxWithdrawable, minWithdrawable })}</OptionalLabel>
                    )
                  }
                  placeholder={t('lightning.receive.amountSats.placeholder')}
                  id="amountSatsInput"
                  onInput={onAmountSatsChange}
                  value={amountSatsText}
                  autoFocus
                />
                <Input
                  label={t('lightning.receive.description.label')}
                  placeholder={t('lightning.receive.description.placeholder')}
                  id="descriptionInput"
                  onInput={onDescriptionChange}
                  value={description}
                  labelSection={<OptionalLabel>{t('lightning.receive.description.optional')}</OptionalLabel>}
                />
                <Status hidden={!showOpenChannelWarning} type="info">
                  {t('lightning.receive.openChannelWarning', { feeSat: toSat(openChannelFeeResponse?.feeMsat!) })}
                </Status>
                <Button transparent onClick={onScanQrCode}>
                  {t('lightning.receive.qrCode.label')}
                </Button>
              </Column>
            </Grid>
          </ViewContent>
          <ViewButtons>
            <Button primary onClick={receivePayment} disabled={receiveDisabled}>
              {parsedInput ? t('button.receive') : t('lightning.receive.invoice.create')}
            </Button>
            <Button secondary onClick={back}>
              {t('button.back')}
            </Button>
          </ViewButtons>
        </View>
      );
    case 'scan':
      return (
        <View fitContent>
          <ViewContent textAlign="center">
            <Grid col="1">
              <Column>
                {/* this flickers quickly, as there is 'SdkError: Generic: Breez SDK error: Unrecognized input type' when logging rawInputError */}
                {rawInputError && <Status type="warning">{rawInputError}</Status>}
                <ScanQRVideo onResult={parseInput} />
                {/* Note: unfortunatelly we probably can't read from HTML5 clipboard api directly in Qt/Andoird WebView */}
                <Button transparent onClick={() => console.log('TODO: implement paste')}>
                  {t('lightning.receive.rawInput.label')}
                </Button>
              </Column>
            </Grid>
          </ViewContent>
          <ViewButtons reverseRow>
            {/* <Button primary onClick={parseInput} disabled={busy}>
          {t('button.send')}
        </Button> */}
            <Button secondary onClick={back}>
              {t('button.back')}
            </Button>
          </ViewButtons>
        </View>
      );
    case 'wait':
      return parsedInput ? (
        <Spinner text={t('lightning.receive.receiving.message')} guideExists={false} />
      ) : (
        <Spinner text={t('lightning.receive.invoice.creating')} guideExists={false} />
      );
    case 'invoice':
      return (
        <View fitContent minHeight="100%">
          <ViewContent textAlign="center">
            <Grid col="1">
              <Column>
                <h1 className={styles.title}>{t('lightning.receive.invoice.title')}</h1>
                <div>
                  <QRCode data={receivePaymentResponse?.lnInvoice.bolt11} />
                </div>
                <div className={styles.invoiceSummary}>
                  {amountSatsText} sats (--- EUR)
                  {description && ` / ${description}`}
                </div>
                <ColumnButtons>
                  <CopyButton data={receivePaymentResponse?.lnInvoice.bolt11} successText={t('lightning.receive.invoice.copied')}>
                    {t('button.copy')}
                  </CopyButton>
                  <Button transparent onClick={back}>
                    <EditActive className={styles.btnIcon} />
                    {t('lightning.receive.invoice.edit')}
                  </Button>
                </ColumnButtons>
              </Column>
            </Grid>
          </ViewContent>
          <ViewButtons>
            <Button secondary onClick={() => route('/lightning')}>
              {t('button.done')}
            </Button>
          </ViewButtons>
        </View>
      );
    case 'success':
      return (
        <View fitContent textCenter verticallyCentered>
          <ViewContent withIcon="success">
            <p>{t('lightning.receive.success.message')}</p>
            {amountSatsText} sats (--- EUR)
            <br />
            {description && ` / ${description}`}
          </ViewContent>
        </View>
      );
    }
  };

  return (
    <GuideWrapper>
      <GuidedContent>
        <Main>
          <Status type="warning" hidden={!receiveError}>
            {receiveError}
          </Status>
          <Header title={<h2>{t('lightning.receive.title')}</h2>} />
          {renderSteps()}
        </Main>
      </GuidedContent>
    </GuideWrapper>
  );
}

type TCopyButtonProps = {
  data?: string;
  successText?: string;
  children: string;
};

const CopyButton = ({ data, successText, children }: TCopyButtonProps) => {
  const [state, setState] = useState('ready');
  const [buttonText, setButtonText] = useState(children);

  const copy = () => {
    try {
      if (data) {
        navigator.clipboard.writeText(data).then(() => {
          setState('success');
          successText && setButtonText(successText);
        });
      }
    } catch (error) {
      setState('ready');
      if (error instanceof Error) {
        setButtonText(error.message);
      }
      setButtonText(`${error}`);
    }
  };

  return (
    <Button transparent onClick={copy} disabled={!data}>
      {state === 'success' ? <Checked className={styles.btnIcon} /> : <Copy className={styles.btnIcon} />}
      {buttonText}
    </Button>
  );
};
