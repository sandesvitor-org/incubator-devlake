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

import { useState, useEffect, useMemo } from 'react';

export interface UseRowSelectionProps<T> {
  dataSource: T[];
  rowSelection?: {
    rowKey?: ID;
    getRowKey?: (data: T) => ID;
    type?: 'checkbox' | 'radio';
    selectedRowKeys?: ID[];
    onChange?: (selectedRowKeys: ID[]) => void;
  };
}

export const useRowSelection = <T>({ dataSource, rowSelection }: UseRowSelectionProps<T>) => {
  const [selectedKeys, setSelectedKeys] = useState<ID[]>([]);

  const {
    rowKey = 'id',
    getRowKey = (data: T) => (data as any)[rowKey],
    type = 'checkbox',
    selectedRowKeys,
    onChange,
  } = {
    rowKey: 'id',
    type: 'checkbox',
    ...rowSelection,
  };

  useEffect(() => {
    setSelectedKeys(selectedRowKeys ?? []);
  }, [selectedRowKeys]);

  const handleChecked = (data: T) => {
    const key = getRowKey(data);
    let result: ID[] = selectedKeys;

    switch (true) {
      case !selectedKeys.includes(key) && type === 'radio':
        result = [key];
        break;
      case !selectedKeys.includes(key) && type === 'checkbox':
        result = [...selectedKeys, key];
        break;
      case selectedKeys.includes(key) && type === 'checkbox':
        result = selectedKeys.filter((k) => k !== key);
        break;
    }

    onChange ? onChange(result) : setSelectedKeys(result);
  };

  const handleCheckedAll = () => {
    let result: string[] = [];

    if (selectedKeys.length !== dataSource.length) {
      result = dataSource.map(getRowKey);
    }

    onChange ? onChange(result) : setSelectedKeys(result);
  };

  return useMemo(
    () => ({
      canSelection: !!rowSelection,
      selectionType: type,
      getCheckedAll: () => dataSource.length === selectedKeys.length,
      onCheckedAll: handleCheckedAll,
      getChecked: (data: T) => {
        return selectedKeys.includes(getRowKey(data));
      },
      onChecked: handleChecked,
    }),
    [selectedKeys],
  );
};
