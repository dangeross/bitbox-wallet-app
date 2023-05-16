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

import { apiGet, apiPost } from './request';

type TConfig = {
  backend?: unknown;
  frontend?: unknown;
};

let pendingConfig: TConfig = {};

/**
 * expects an object with a backend or frontend key
 * i.e. { frontend: { language }}
 * returns a promise and passes the new config
 */
export const setConfig = (object: TConfig) => {
  return apiGet('config')
    .then((currentConfig = {}) => {
      const nextConfig = Object.assign(currentConfig, {
        backend: Object.assign({}, currentConfig.backend, pendingConfig.backend, object.backend),
        frontend: Object.assign({}, currentConfig.frontend, pendingConfig.frontend, object.frontend)
      });
      pendingConfig = nextConfig;
      return apiPost('config', nextConfig)
        .then(() => {
          pendingConfig = {};
          return nextConfig;
        });
    });
};