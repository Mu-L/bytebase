import { create } from "@bufbuild/protobuf";
import { createContextValues } from "@connectrpc/connect";
import { settingServiceClientConnect } from "@/grpcweb";
import { silentContextKey } from "@/grpcweb/context-key";
import { 
  GetSettingRequestSchema, 
  UpdateSettingRequestSchema
} from "@/types/proto-es/v1/setting_service_pb";
import { convertNewSettingToOld, convertOldSettingToNew, convertOldSettingNameToNew } from "@/utils/v1/setting-conversions";
import type { ResourceId } from "@/types";
import { unknownEnvironment } from "@/types";
import {
  EnvironmentSetting,
  EnvironmentSetting_Environment,
  Setting_SettingName as OldSettingName,
} from "@/types/proto/v1/setting_service";
import type { Environment } from "@/types/v1/environment";
import { orderBy } from "lodash-es";
import { defineStore } from "pinia";
import { computed } from "vue";
import { environmentNamePrefix } from "./common";

interface EnvironmentState {
  environmentMapById: Map<ResourceId, Environment>;
}

const getEnvironmentByIdMap = (
  environments: Environment[]
): Map<ResourceId, Environment> => {
  return new Map(
    environments.map((environment) => [environment.id, environment])
  );
};

const convertToEnvironments = (
  environments: EnvironmentSetting_Environment[]
): Environment[] => {
  return environments.map<Environment>((env, i) => {
    return {
      name: `${environmentNamePrefix}${env.id}`,
      id: env.id,
      title: env.title,
      order: i,
      color: env.color,
      tags: env.tags,
    };
  });
};

const convertEnvironment = (
  env: Environment
): EnvironmentSetting_Environment => {
  const res: EnvironmentSetting_Environment = {
    name: env.name,
    id: env.id,
    title: env.title,
    color: env.color,
    tags: env.tags,
  };
  return res;
};

const convertEnvironments = (
  environments: Environment[]
): EnvironmentSetting_Environment[] => {
  return environments.map(convertEnvironment);
};

const getEnvironmentSetting = async (
  silent = false
): Promise<Environment[]> => {
  const newName = convertOldSettingNameToNew(OldSettingName.ENVIRONMENT);
  const request = create(GetSettingRequestSchema, {
    name: `settings/${newName}`,
  });
  const response = await settingServiceClientConnect.getSetting(request, {
    contextValues: createContextValues().set(silentContextKey, silent),
  });
  // Extract environments from new proto format
  if (response.value?.value?.case === "environmentSetting") {
    const oldSetting = convertNewSettingToOld(response);
    const settingEnvironments =
      oldSetting.value?.environmentSetting?.environments ?? [];
    return convertToEnvironments(settingEnvironments);
  }
  return [];
};

const updateEnvironmentSetting = async (
  environment: EnvironmentSetting
): Promise<Environment[]> => {
  // Create old setting object and convert to new format
  const newName = convertOldSettingNameToNew(OldSettingName.ENVIRONMENT);
  const oldSetting = {
    name: `settings/${newName}`,
    value: {
      environmentSetting: environment,
    },
  };
  const newSetting = convertOldSettingToNew(oldSetting);
  
  const request = create(UpdateSettingRequestSchema, {
    setting: newSetting,
    updateMask: { paths: ["environment_setting"] },
  });
  const response = await settingServiceClientConnect.updateSetting(request);
  
  // Extract environments from response
  if (response.value?.value?.case === "environmentSetting") {
    const oldResponse = convertNewSettingToOld(response);
    const settingEnvironments =
      oldResponse.value?.environmentSetting?.environments ?? [];
    return convertToEnvironments(settingEnvironments);
  }
  return [];
};

export const useEnvironmentV1Store = defineStore("environment_v1", {
  state: (): EnvironmentState => ({
    environmentMapById: new Map(),
  }),
  getters: {
    environmentList(state) {
      return orderBy(
        Array.from(state.environmentMapById.values()),
        (env) => env.order,
        "asc"
      );
    },
  },
  actions: {
    async fetchEnvironments(silent = false) {
      const environments = await getEnvironmentSetting(silent);
      this.environmentMapById = getEnvironmentByIdMap(environments);
      return environments;
    },
    getEnvironmentList(): Environment[] {
      return this.environmentList;
    },
    async createEnvironment(
      environment: Partial<Environment>
    ): Promise<Environment> {
      const e: EnvironmentSetting_Environment = {
        name: "",
        id: environment.id ?? "",
        title: environment.title ?? "",
        color: environment.color ?? "",
        tags: environment.tags ?? {},
      };
      const newEnvironmentSettingValue = [
        ...convertEnvironments(this.environmentList),
        e,
      ];

      const newEnvironments = await updateEnvironmentSetting({
        environments: newEnvironmentSettingValue,
      });

      const newEnvironmentMapById = getEnvironmentByIdMap(newEnvironments);
      this.environmentMapById = newEnvironmentMapById;
      const newEnvironment = newEnvironmentMapById.get(environment.id ?? "");
      if (!newEnvironment) {
        throw new Error(`environment with id ${environment.id} not found`);
      }
      return newEnvironment;
    },
    async updateEnvironment(
      update: Partial<Environment>
    ): Promise<Environment> {
      const originData = await this.getOrFetchEnvironmentByName(
        update.id || ""
      );
      if (!originData) {
        throw new Error(`environment with id ${update.id} not found`);
      }
      const newEnvironments = await updateEnvironmentSetting({
        environments: convertEnvironments(
          this.environmentList.map((environment) => {
            if (environment.id === update.id) {
              environment.title = update.title ?? environment.title;
              environment.color = update.color ?? environment.color;
              environment.tags = update.tags ?? environment.tags;
              environment.order = update.order ?? environment.order;
            }
            return environment;
          })
        ),
      });

      const newEnvironmentMapById = getEnvironmentByIdMap(newEnvironments);
      this.environmentMapById = newEnvironmentMapById;
      const newEnvironment = newEnvironmentMapById.get(update.id ?? "");
      if (!newEnvironment) {
        throw new Error(`environment with id ${update.id} not found`);
      }
      return newEnvironment;
    },
    async deleteEnvironment(name: string): Promise<void> {
      const id = name.replace(environmentNamePrefix, "");
      const newEnvironments = await updateEnvironmentSetting({
        environments: convertEnvironments(
          this.environmentList.filter((environment) => environment.id !== id)
        ),
      });
      this.environmentMapById = getEnvironmentByIdMap(newEnvironments);
    },
    async reorderEnvironmentList(
      orderedEnvironmentList: Environment[]
    ): Promise<Environment[]> {
      const newEnvironments = await updateEnvironmentSetting({
        environments: convertEnvironments(orderedEnvironmentList),
      });
      this.environmentMapById = getEnvironmentByIdMap(newEnvironments);
      return newEnvironments;
    },
    async getOrFetchEnvironmentByName(
      name: string,
      silent = false
    ): Promise<Environment | undefined> {
      const id = name.replace(environmentNamePrefix, "");
      const cachedData = this.environmentMapById.get(id);
      if (cachedData) {
        return cachedData;
      }
      await this.fetchEnvironments(silent);
      const environment = this.environmentMapById.get(id);
      return environment;
    },
    getEnvironmentByName(name: string) {
      const id = name.replace(environmentNamePrefix, "");
      return this.environmentMapById.get(id) ?? unknownEnvironment();
    },
  },
});

export const useEnvironmentV1List = () => {
  const store = useEnvironmentV1Store();
  return computed(() => store.getEnvironmentList());
};
