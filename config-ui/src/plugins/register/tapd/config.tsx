/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import { ExternalLink } from '@/components';
import { DOC_URL } from '@/release';

import type { PluginConfigType } from '../../types';
import { PluginType } from '../../types';

import { DataScope } from './data-scope';
import Icon from './assets/icon.svg';

export const TAPDConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'tapd',
  name: 'TAPD',
  icon: Icon,
  sort: 9,
  connection: {
    docLink: DOC_URL.PLUGIN.TAPD.BASIS,
    initialValues: {
      endpoint: 'https://api.tapd.cn',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel: 'You do not need to enter the endpoint URL, because all versions use the same URL.',
        disabled: true,
      },
      {
        key: 'username',
        label: 'API Account',
        subLabel: (
          <span>
            Please follow the instruction <ExternalLink link={DOC_URL.PLUGIN.TAPD.USERNAMEPASSWORD}>here</ExternalLink>{' '}
            to find your API account information.
          </span>
        ),
        placeholder: 'Your API Account',
      },
      {
        key: 'password',
        label: 'API Token',
        placeholder: 'Your API Token',
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 3,000 requests/hour for data collection for TAPD. But you can adjust the collection speed by setting up your desirable rate limit.',
        learnMore: DOC_URL.PLUGIN.TAPD.RATE_LIMIT,
        externalInfo: 'The maximum rate limit of TAPD is 3,600 requests/hour.',
        defaultValue: 3000,
      },
    ],
  },
  dataScope: {
    render: ({ ...props }) => <DataScope {...props} />,
  },
  scopeConfig: {
    entities: ['TICKET', 'CROSS'],
    transformation: {
      typeMappings: {},
      statusMappings: {},
    },
  },
};
