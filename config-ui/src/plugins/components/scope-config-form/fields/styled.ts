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

export const AdditionalSettings = styled.div`
  h2 {
    display: flex;
    align-items: center;

    .bp4-switch {
      margin: 0;
    }
  }

  .radio {
    display: flex;
    align-items: center;
    margin: 8px 0 16px;

    p {
      margin: 0;
    }

    .bp4-control {
      margin: 0;
    }
  }

  .refdiff {
    display: flex;
    align-items: center;
    padding-left: 20px;

    .bp4-input-group {
      margin: 0 8px;
    }
  }
`;
