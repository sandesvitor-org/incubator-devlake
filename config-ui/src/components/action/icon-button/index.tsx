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

import { Button, Intent, Position, IconName } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';

interface Props {
  style?: React.CSSProperties;
  icon?: IconName;
  image?: JSX.Element;
  tooltip: string;
  loading?: boolean;
  disabled?: boolean;
  onClick?: () => void;
}

export const IconButton = ({ style, icon, image, tooltip, loading, disabled, onClick }: Props) => {
  return (
    <Tooltip2 intent={Intent.PRIMARY} position={Position.TOP} content={tooltip}>
      <Button
        style={style}
        loading={loading}
        disabled={disabled}
        minimal
        intent={Intent.PRIMARY}
        icon={icon ?? image}
        onClick={onClick}
      />
    </Tooltip2>
  );
};
