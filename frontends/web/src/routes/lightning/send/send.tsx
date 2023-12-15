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

import { ChangeEvent, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import * as accountApi from '../../../api/account';
import { Column, Grid, GuideWrapper, GuidedContent, Header, Main } from '../../../components/layout';
import { View, ViewButtons, ViewContent } from '../../../components/view/view';
import { Button, Input, OptionalLabel } from '../../../components/forms';
import {
  InputType,
  InputTypeVariant,
  LnInvoice,
  LnUrlPayRequestData,
  SdkError,
  getParseInput,
  postLnurlPay,
  postSendPayment
} from '../../../api/lightning';
import { SimpleMarkup } from '../../../utils/markup';
import { route } from '../../../utils/route';
import { toMsat, toSat } from '../../../utils/conversion';
import { Amount } from '../../../components/amount/amount';
import { FiatConversion } from '../../../components/rates/rates';
import { Status } from '../../../components/status/status';
import { ScanQRVideo } from '../../account/send/components/inputs/scan-qr-video';
import { Spinner } from '../../../components/spinner/Spinner';
import styles from './send.module.css';

type TStep = 'select-invoice' | 'confirm' | 'sending' | 'success';

const SendingSpinner = () => {
  const { t } = useTranslation();
  // Show dummy connecting-to-server message first
  const [message, setStep] = useState<string>(t('lightning.send.sending.connecting'));

  setTimeout(() => {
    setStep(t('lightning.send.sending.message'));
  }, 4000);

  return <Spinner text={message} guideExists={false} />;
};

export function Send() {
  const { t } = useTranslation();
  const [amountRequired, setAmountRequired] = useState<boolean>(false);
  const [amountSats, setAmountSats] = useState<number>(0);
  const [amountSatsText, setAmountSatsText] = useState<string>('');
  const [description, setDescription] = useState<string>('');
  const [minSendable, setMinSendable] = useState<number>(0);
  const [maxSendable, setMaxSendable] = useState<number>(0);
  const [optionalComment, setOptionalComment] = useState<string>('');
  const [parsedInput, setParsedInput] = useState<InputType>();
  const [rawInputError, setRawInputError] = useState<string>();
  const [sendDisabled, setSendDisabled] = useState<boolean>();
  const [sendError, setSendError] = useState<string>();
  const [step, setStep] = useState<TStep>('select-invoice');

  const back = () => {
    setSendError(undefined);
    switch (step) {
    case 'select-invoice':
      route('/lightning');
      break;
    case 'confirm':
    case 'success':
      setStep('select-invoice');
      setParsedInput(undefined);
      break;
    }
  };

  const onAmountSatsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target as HTMLInputElement;
    setAmountSatsText(target.value);
  };

  const onOptionalCommentChange = (event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target as HTMLInputElement;
    setOptionalComment(target.value);
  };

  useEffect(() => {
    setAmountSats(+amountSatsText);
  }, [amountSatsText]);

  useEffect(() => {
    switch (parsedInput?.type) {
    case InputTypeVariant.BOLT11:
      setSendDisabled(amountRequired && amountSats <= 0);
      break;
    case InputTypeVariant.LN_URL_PAY:
      setSendDisabled(
        amountSats <= 0 ||
          amountSats < minSendable ||
          amountSats > maxSendable ||
          optionalComment.length > parsedInput.data.commentAllowed
      );
      break;
    }
  }, [amountRequired, amountSatsText, amountSats, maxSendable, minSendable, parsedInput, optionalComment]);

  useEffect(() => {
    if (parsedInput?.type === InputTypeVariant.LN_URL_PAY) {
      const metadata: any[] = JSON.parse(parsedInput.data.metadataStr);
      const metadataTextPlain = metadata.find((el) => {
        return el[0] === 'text/plain';
      });
      const metadataTextIdentifier = metadata.find((el) => {
        return el[0] === 'text/identifier';
      });
      setDescription(
        metadataTextPlain
          ? metadataTextPlain[1]
          : t('lightning.send.confirm.lnUrl.title') + metadataTextIdentifier
            ? ` ${metadataTextIdentifier[1]}`
            : ''
      );
    }
  }, [parsedInput, t]);

  const parseInput = useCallback(async (rawInput: string) => {
    setRawInputError(undefined);
    try {
      const result = await getParseInput({ s: rawInput });
      switch (result.type) {
      case InputTypeVariant.BOLT11:
        if (result.invoice.amountMsat) {
          setAmountRequired(false);
          setAmountSatsText(toSat(result.invoice.amountMsat).toString());
        } else {
          setAmountRequired(true);
          setAmountSatsText('');
        }
        setParsedInput(result);
        setStep('confirm');
        break;
      case InputTypeVariant.LN_URL_PAY:
        setAmountRequired(true);
        setAmountSatsText('');
        setMaxSendable(toSat(result.data.maxSendable));
        setMinSendable(toSat(result.data.minSendable));
        setOptionalComment('');
        setParsedInput(result);
        setStep('confirm');
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

  const sendPayment = async () => {
    setStep('sending');
    setSendError(undefined);
    try {
      switch (parsedInput?.type) {
      case InputTypeVariant.BOLT11:
        await postSendPayment({ bolt11: parsedInput.invoice.bolt11, amountMsat: amountRequired ? toMsat(amountSats) : undefined });
        setStep('success');
        setTimeout(() => route('/lightning'), 5000);
        break;
      case InputTypeVariant.LN_URL_PAY:
        await postLnurlPay({ data: parsedInput.data, amountMsat: toMsat(amountSats), comment: optionalComment });
        setStep('success');
        setTimeout(() => route('/lightning'), 5000);
        break;
      }
    } catch (e) {
      setStep('select-invoice');
      if (e instanceof SdkError) {
        setSendError(e.message);
      } else {
        setSendError(String(e));
      }
    }
  };

  const renderInputTypes = () => {
    switch (parsedInput!.type) {
    case InputTypeVariant.BOLT11:
      return renderBolt11Inputs(parsedInput!.invoice);
    case InputTypeVariant.LN_URL_PAY:
      return renderLnUrlPayInputs(parsedInput!.data);
    }
  };

  const renderBolt11Inputs = (invoice: LnInvoice) => {
    const balance: accountApi.IBalance = {
      hasAvailable: true,
      available: {
        amount: `${toSat(invoice.amountMsat || 0)}`,
        unit: 'sat'
      },
      hasIncoming: false,
      incoming: {
        amount: '0',
        unit: 'sat'
      }
    };
    return (
      <Column>
        <h1 className={styles.title}>{t('lightning.send.confirm.title')}</h1>
        {amountRequired ? (
          <Input
            type="number"
            min="0"
            label={t('lightning.send.amountSats.label')}
            placeholder={t('lightning.send.amountSats.placeholder')}
            id="amountSatsInput"
            onInput={onAmountSatsChange}
            value={amountSatsText}
            autoFocus
          />
        ) : (
          <div className={styles.info}>
            <h2 className={styles.label}>{t('lightning.send.confirm.amount')}</h2>
            <Amount amount={balance.available.amount} unit={balance.available.unit} removeBtcTrailingZeroes />/{' '}
            <FiatConversion amount={balance.available} noBtcZeroes />
          </div>
        )}
        {invoice.description && (
          <div className={styles.info}>
            <h2 className={styles.label}>{t('lightning.send.confirm.memo')}</h2>
            {invoice.description}
          </div>
        )}
      </Column>
    );
  };

  const renderLnUrlPayInputs = (data: LnUrlPayRequestData) => {
    return (
      <Column>
        <h1 className={styles.title}>{ description }</h1>
        <Input
          type="number"
          min="0"
          label={t('lightning.send.amountSats.label')}
          labelSection={<OptionalLabel>{t('lightning.send.amountSats.limitLabel', { maxSendable, minSendable })}</OptionalLabel>}
          placeholder={t('lightning.send.amountSats.placeholder')}
          id="amountSatsInput"
          onInput={onAmountSatsChange}
          value={amountSatsText}
          autoFocus
        />
        {data.commentAllowed > 0 && (
          <Input
            type="text"
            label={t('lightning.send.comment.label')}
            labelSection={<OptionalLabel>{t('lightning.send.comment.optional')}</OptionalLabel>}
            placeholder={t('lightning.send.comment.placeholder')}
            id="optionalComment"
            onInput={onOptionalCommentChange}
            value={optionalComment}
          />
        )}
      </Column>
    );
  };

  const renderSteps = () => {
    switch (step) {
    case 'select-invoice':
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
                  {t('lightning.send.rawInput.label')}
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
    case 'confirm':
      return (
        <View fitContent minHeight="100%">
          <ViewContent>
            <Grid col="1">{renderInputTypes()}</Grid>
          </ViewContent>
          <ViewButtons>
            <Button primary onClick={sendPayment} disabled={sendDisabled}>
              {t('button.send')}
            </Button>
            <Button secondary onClick={back}>
              {t('button.back')}
            </Button>
          </ViewButtons>
        </View>
      );
    case 'sending':
      return <SendingSpinner />;
    case 'success':
      return (
        <View fitContent textCenter verticallyCentered>
          <ViewContent withIcon="success">
            <SimpleMarkup className={styles.successMessage} markup={t('lightning.send.success.message')} tagName="p" />
          </ViewContent>
        </View>
      );
    }
  };

  return (
    <GuideWrapper>
      <GuidedContent>
        <Main>
          <Status type="warning" hidden={!sendError}>
            {sendError}
          </Status>
          <Header title={<h2>{t('lightning.send.title')}</h2>} />
          {renderSteps()}
        </Main>
      </GuidedContent>
    </GuideWrapper>
  );
}
