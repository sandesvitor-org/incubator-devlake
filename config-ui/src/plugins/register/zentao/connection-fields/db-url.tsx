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

import { useEffect } from 'react';
import { InputGroup } from '@blueprintjs/core';

import { FormItem } from '@/components';

interface Props {
  initialValue: string;
  value: string;
  setValue: (value: string) => void;
}

export const DBUrl = ({ initialValue, value, setValue }: Props) => {
  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  return (
    <FormItem label="Database URL" subLabel="Provide the DB URL of Zentao if you want to collect issue changelogs.">
      <InputGroup
        placeholder="e.g. mysql://root:devlake@sshd-proxy:3306/zentao?charset=utf8mb4&parseTime=True"
        value={value}
        onChange={handleChange}
      />
    </FormItem>
  );
};
