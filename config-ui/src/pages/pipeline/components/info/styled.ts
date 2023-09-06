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

import { Colors } from '@blueprintjs/core';
import styled from 'styled-components';

import { Card } from '@/components';

export const Wrapper = styled(Card)`
  ul {
    display: flex;
    align-items: center;
  }

  li {
    flex: 5;
    display: flex;
    flex-direction: column;

    &:last-child {
      flex: 1;
    }

    & > span {
      font-size: 12px;
      color: #94959f;
      text-align: center;
    }

    & > strong {
      display: flex;
      align-items: center;
      justify-content: center;
      margin-top: 8px;
    }
  }

  p.message {
    margin: 8px 0 0;
    color: ${Colors.RED3};
  }
`;
