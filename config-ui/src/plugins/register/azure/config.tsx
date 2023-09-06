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

import Icon from './assets/icon.svg';
import { BaseURL } from './connection-fields';

export const AzureConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'azuredevops',
  name: 'Azure DevOps',
  icon: Icon,
  sort: 6,
  connection: {
    docLink: DOC_URL.PLUGIN.AZUREDEVOPS.BASIS,
    fields: [
      'name',
      () => <BaseURL key="base-url" />,
      {
        key: 'token',
        label: 'Personal Access Token',
        subLabel: (
          <span>
            <ExternalLink link={DOC_URL.PLUGIN.AZUREDEVOPS.AUTH_TOKEN}>Learn about how to create a PAT</ExternalLink>{' '}
            Please select ALL ACCESSIBLE ORGANIZATIONS for the Organization field when you create the PAT.
          </span>
        ),
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 18,000 requests/hour for data collection for Azure DevOps. But you can adjust the collection speed by setting up your desirable rate limit.',
        learnMore: DOC_URL.PLUGIN.AZUREDEVOPS.RATE_LIMIT,
        externalInfo: 'Azure DevOps does not specify a maximum value of rate limit.',
        maximum: 18000,
      },
    ],
  },
  dataScope: {
    millerColumns: {
      title: 'Add Repositories by Selecting from the Directory',
      subTitle: 'The following directory lists out all repositories in your organizations.',
    },
  },
  scopeConfig: {
    entities: ['CODE', 'CODEREVIEW', 'CROSS', 'CICD'],
    transformation: {
      deploymentPattern: '(deploy|push-image)',
      productionPattern: 'prod(.*)',
      refdiff: {
        tagsLimit: 10,
        tagsPattern: '/v\\d+\\.\\d+(\\.\\d+(-rc)*\\d*)*$/',
      },
    },
  },
};
