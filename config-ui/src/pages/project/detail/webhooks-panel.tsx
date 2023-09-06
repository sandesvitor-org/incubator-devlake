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

import { useState, useMemo } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { Alert, NoData } from '@/components';
import type { WebhookItemType } from '@/plugins/register/webook';
import { WebhookCreateDialog, WebhookSelectorDialog, WebHookConnection } from '@/plugins/register/webook';
import { operator } from '@/utils';

import type { ProjectType } from '../types';

import * as API from './api';

interface Props {
  project: ProjectType;
  onRefresh: () => void;
}

export const WebhooksPanel = ({ project, onRefresh }: Props) => {
  const [type, setType] = useState<'selectExist' | 'create'>();
  const [operating, setOperating] = useState(false);

  const webhookIds = useMemo(
    () =>
      project.blueprint
        ? project.blueprint.settings.connections
            .filter((cs: any) => cs.plugin === 'webhook')
            .map((cs: any) => cs.connectionId)
        : [],
    [project],
  );

  const handleCancel = () => {
    setType(undefined);
  };

  const handleCreate = async (id: ID) => {
    const payload = {
      ...project.blueprint,
      settings: {
        ...project.blueprint.settings,
        connections: [
          ...project.blueprint.settings.connections,
          {
            plugin: 'webhook',
            connectionId: id,
          },
        ],
      },
    };

    const [success] = await operator(() => API.updateBlueprint(project.blueprint.id, payload), {
      setOperating,
    });

    if (success) {
      onRefresh();
    }
  };

  const handleSelect = async (items: WebhookItemType[]) => {
    const payload = {
      ...project.blueprint,
      settings: {
        ...project.blueprint.settings,
        connections: [
          ...project.blueprint.settings.connections,
          ...items.map((it) => ({
            plugin: 'webhook',
            connectionId: it.id,
          })),
        ],
      },
    };

    const [success] = await operator(() => API.updateBlueprint(project.blueprint.id, payload), {
      setOperating,
    });

    if (success) {
      onRefresh();
    }
  };

  const handleDelete = async (id: ID) => {
    const payload = {
      ...project.blueprint,
      settings: {
        ...project.blueprint.settings,
        connections: project.blueprint.settings.connections.filter(
          (cs: any) => !(cs.plugin === 'webhook' && cs.connectionId === id),
        ),
      },
    };

    const [success] = await operator(() => API.updateBlueprint(project.blueprint.id, payload), {
      setOperating,
    });

    if (success) {
      onRefresh();
    }
  };

  return (
    <>
      <Alert
        style={{ marginBottom: 24, color: '#3C5088' }}
        content={
          <>
            <div>
              The data pushed by Webhooks will only be calculated for DORA in the next run of the Blueprint of this
              project because DORA relies on the post-processing of "deployments," "incidents," and "pull requests"
              triggered by running the blueprint.
            </div>
            <div style={{ marginTop: 16 }}>
              To calculate DORA after receiving Webhook data immediately, you can visit the{' '}
              <b style={{ textDecoration: 'underline' }}>Status tab</b> of the Blueprint page and click on Run Now.
            </div>
          </>
        }
      />
      {!webhookIds.length ? (
        <>
          <NoData
            text="Push `incidents` or `deployments` from your tools by incoming webhooks."
            action={
              <>
                <Button intent={Intent.PRIMARY} icon="plus" text="Add a Webhook" onClick={() => setType('create')} />
                <div style={{ margin: '8px 0' }}>or</div>
                <Button
                  outlined
                  intent={Intent.PRIMARY}
                  text="Select Existing Webhooks"
                  onClick={() => setType('selectExist')}
                />
              </>
            }
          />
          {type === 'create' && <WebhookCreateDialog isOpen onCancel={handleCancel} onSubmitAfter={handleCreate} />}
          {type === 'selectExist' && (
            <WebhookSelectorDialog isOpen saving={operating} onCancel={handleCancel} onSubmit={handleSelect} />
          )}
        </>
      ) : (
        <WebHookConnection filterIds={webhookIds} onCreateAfter={handleCreate} onDeleteAfter={handleDelete} />
      )}
    </>
  );
};
