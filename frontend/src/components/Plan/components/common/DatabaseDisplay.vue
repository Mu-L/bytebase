<template>
  <div :class="containerClass">
    <InstanceV1EngineIcon
      v-if="instanceResource"
      class="inline-block mr-1"
      :instance="instanceResource"
      :size="props.size || 'small'"
    />
    <span v-if="showEnvironment && environment" class="text-gray-500 mr-1">
      ({{ environment.title }})
    </span>
    <span class="truncate text-gray-600">
      {{ instanceDisplayName }}
    </span>
    <ChevronRightIcon :class="chevronClass" />
    <span class="truncate text-gray-800">
      {{ databaseDisplayName }}
    </span>
    <template v-if="project">
      <router-link
        class="pl-2"
        :to="databaseV1UrlWithProject(project, database)"
        target="_blank"
        @click.stop
      >
        <ExternalLinkIcon :size="16" />
      </router-link>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ChevronRightIcon, ExternalLinkIcon } from "lucide-vue-next";
import { computed } from "vue";
import { InstanceV1EngineIcon } from "@/components/v2";
import {
  useDatabaseV1Store,
  useEnvironmentV1Store,
  useInstanceV1Store,
} from "@/store";
import { isValidDatabaseName, unknownInstance } from "@/types";
import {
  databaseV1UrlWithProject,
  extractDatabaseResourceName,
  extractInstanceResourceName,
} from "@/utils";

const props = defineProps<{
  database: string;
  project?: string;
  showEnvironment?: boolean;
  size?: "small" | "medium" | "large";
}>();

const environmentStore = useEnvironmentV1Store();
const instanceStore = useInstanceV1Store();

const databaseEntity = computed(() =>
  useDatabaseV1Store().getDatabaseByName(props.database)
);

const instanceResource = computed(() => {
  if (isValidDatabaseName(databaseEntity.value.name)) {
    return databaseEntity.value.instanceResource;
  }
  // Extract instance name from the database name (which is actually a target string)
  const instanceName = extractInstanceResourceName(props.database);
  if (instanceName) {
    const instance = instanceStore.getInstanceByName(
      `instances/${instanceName}`
    );
    if (instance.name !== unknownInstance().name) {
      return instance;
    }
  }

  return null; // Don't show instance icon if we can't resolve it
});

const instanceDisplayName = computed(() => {
  if (instanceResource.value) {
    return instanceResource.value.title;
  }
  return extractInstanceResourceName(props.database) || "Unknown Instance";
});

const databaseDisplayName = computed(() => {
  const { databaseName } = extractDatabaseResourceName(props.database);
  return databaseName || "Unknown Database";
});

const environment = computed(() => {
  const environmentName =
    databaseEntity.value.environment || instanceResource.value?.environment;
  if (!environmentName) {
    return undefined;
  }

  return environmentStore.getEnvironmentByName(environmentName);
});

const containerClass = computed(() => {
  const baseClass = "flex items-center truncate";
  const sizeClasses = {
    small: "text-xs",
    medium: "text-sm",
    large: "text-base",
  };
  return [baseClass, sizeClasses[props.size || "medium"]];
});

const chevronClass = computed(() => {
  const baseClass = "inline opacity-60 text-gray-500 shrink-0";
  const sizeClasses = {
    small: "w-3 h-3",
    medium: "w-4 h-4",
    large: "w-5 h-5",
  };
  return [baseClass, sizeClasses[props.size || "medium"]];
});
</script>
