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

import { DOC_URL } from '@/release';

import type { PluginConfigType } from '../../types';
import { PluginType } from '../../types';

import Icon from './assets/icon.png';

export const PokemonConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'pokemon',
  name: 'Pokemon',
  icon: Icon,
  sort: 8,
  connection: {
    docLink: DOC_URL.PLUGIN.PAGERDUTY.BASIS,
    initialValues: {
      endpoint: 'https://pokeapi.co/api/v2/',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        multipleVersions: {
          cloud: 'https://pokeapi.co/api/v2/',
          server: '',
        },
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit for optimized data collection for Bamboo. But you can adjust the collection speed by entering a fixed value. Please note: the rate limit setting applies to all tokens you have entered above.',
        learnMore: DOC_URL.PLUGIN.BAMBOO.RATE_LIMIT,
        externalInfo: 'Pokemon does not specify a maximum value of rate limit.',
        defaultValue: 10000,
      },
    ],
  },
  dataScope: {
    millerColumns: {
      title: 'Pokemon types *',
      subTitle: 'You can either add types by searching or selecting from the following directory.',
    },
  },
};
