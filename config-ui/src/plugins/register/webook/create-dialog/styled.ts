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

import styled from 'styled-components';
import { Colors } from '@blueprintjs/core';

export const Wrapper = styled.div`
  h2 {
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0;
    margin-bottom: 16px;
    padding: 0;
    font-weight: 600;
    color: ${Colors.GREEN5};

    .bp4-icon {
      margin-right: 8px;
    }
  }

  h3 {
    margin: 0;
    padding: 0;
  }

  .block {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
    padding: 10px 16px;
    background: #f0f4fe;
  }
`;
