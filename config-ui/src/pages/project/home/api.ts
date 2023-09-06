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

import { request } from '@/utils';

import type { BlueprintType } from '@/pages/blueprint';
import type { PipelineType } from '@/pages/pipeline';

type GetProjectsParams = {
  page: number;
  pageSize: number;
};

type GetProjectsResponse = {
  projects: Array<{
    name: string;
    createdAt: string;
    blueprint: BlueprintType;
    lastPipeline: PipelineType;
  }>;
  count: number;
};

export const getProjects = (params: GetProjectsParams): Promise<GetProjectsResponse> =>
  request('/projects', { data: params });

type CreateProjectPayload = {
  name: string;
  description: string;
  metrics: Array<{
    pluginName: string;
    pluginOption: string;
    enable: boolean;
  }>;
};

export const createProject = (payload: CreateProjectPayload) =>
  request('/projects', {
    method: 'post',
    data: payload,
  });

export const createBlueprint = (payload: any) =>
  request('/blueprints', {
    method: 'post',
    data: payload,
  });
