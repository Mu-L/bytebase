<template>
  <NDataTable
    key="instance-table"
    size="small"
    :columns="columnList"
    :data="instanceList"
    :striped="true"
    :bordered="bordered"
    :loading="loading"
    :row-key="(data: ComposedInstance) => data.name"
    :checked-row-keys="selectedInstanceNames"
    :row-props="rowProps"
    :paginate-single-page="false"
    @update:checked-row-keys="
        (val) => $emit('update:selected-instance-names', val as string[])
      "
  />
</template>

<script setup lang="tsx">
import {
  ExternalLinkIcon,
  ChevronDownIcon,
  ChevronUpIcon,
} from "lucide-vue-next";
import { NButton, NDataTable, type DataTableColumn } from "naive-ui";
import { computed, reactive, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { InstanceV1Name } from "@/components/v2";
import type { ComposedInstance } from "@/types";
import { urlfy, hostPortOfInstanceV1, hostPortOfDataSource } from "@/utils";
import EnvironmentV1Name from "../../EnvironmentV1Name.vue";

type InstanceDataTableColumn = DataTableColumn<ComposedInstance> & {
  hide?: boolean;
};

interface LocalState {
  dataSourceToggle: Set<string>;
  processing: boolean;
}

const props = withDefaults(
  defineProps<{
    instanceList: ComposedInstance[];
    bordered?: boolean;
    loading?: boolean;
    showSelection?: boolean;
    defaultExpandDataSource?: string[];
    selectedInstanceNames?: string[];
    disabledSelectedInstanceNames?: Set<string>;
    showAddress?: boolean;
    showExternalLink?: boolean;
    onClick?: (instance: ComposedInstance, e: MouseEvent) => void;
  }>(),
  {
    bordered: true,
    showSelection: true,
    onClick: undefined,
    showAddress: true,
    showExternalLink: true,
    disabledSelectedInstanceNames: () => new Set(),
    defaultExpandDataSource: () => [],
    selectedInstanceNames: () => [],
  }
);

defineEmits<{
  (event: "update:selected-instance-names", val: string[]): void;
}>();

const { t } = useI18n();
const router = useRouter();
const state = reactive<LocalState>({
  dataSourceToggle: new Set(),
  processing: false,
});

watch(
  () => props.defaultExpandDataSource,
  () => {
    state.dataSourceToggle = new Set(props.defaultExpandDataSource);
  },
  {
    immediate: true,
    deep: true,
  }
);

const columnList = computed((): InstanceDataTableColumn[] => {
  const SELECTION: InstanceDataTableColumn = {
    type: "selection",
    hide: !props.showSelection,
    disabled: (instance) =>
      props.disabledSelectedInstanceNames.has(instance.name),
    cellProps: () => {
      return {
        onClick: (e: MouseEvent) => {
          e.stopPropagation();
        },
      };
    },
  };
  const NAME: InstanceDataTableColumn = {
    key: "title",
    title: t("common.name"),
    resizable: true,
    render: (instance) => {
      return <InstanceV1Name instance={instance} link={false} />;
    },
  };
  const ENVIRONMENT: InstanceDataTableColumn = {
    key: "environment",
    title: t("common.environment"),
    className: "whitespace-nowrap",
    resizable: true,
    render: (instance) => (
      <EnvironmentV1Name
        environment={instance.environmentEntity}
        link={false}
      />
    ),
  };
  const ADDRESS: InstanceDataTableColumn = {
    key: "address",
    title: t("common.address"),
    resizable: true,
    hide: !props.showAddress,
    render: (instance) => {
      return (
        <div class={"flex items-start gap-x-2"}>
          <div>
            {state.dataSourceToggle.has(instance.name)
              ? instance.dataSources.map((ds) => (
                  <div>{hostPortOfDataSource(ds)}</div>
                ))
              : hostPortOfInstanceV1(instance)}
          </div>
          {instance.dataSources.length > 1 ? (
            <NButton
              quaternary
              size="tiny"
              onClick={(e) => handleDataSourceToggle(e, instance)}
            >
              {state.dataSourceToggle.has(instance.name) ? (
                <ChevronUpIcon class={"w-4 cursor-pointer"} />
              ) : (
                <ChevronDownIcon class={"w-4 cursor-pointer"} />
              )}
            </NButton>
          ) : null}
        </div>
      );
    },
  };
  const EXTERNAL_LINK: InstanceDataTableColumn = {
    key: "project",
    title: t("instance.external-link"),
    resizable: true,
    width: 150,
    hide: !props.showExternalLink,
    render: (instance) =>
      instance.externalLink?.trim().length !== 0 && (
        <NButton
          quaternary
          size="small"
          onClick={() => window.open(urlfy(instance.externalLink), "_blank")}
        >
          <ExternalLinkIcon class="w-4 h-4" />
        </NButton>
      ),
  };
  const LICENSE: InstanceDataTableColumn = {
    key: "instance",
    title: t("subscription.instance-assignment.license"),
    resizable: true,
    width: 100,
    render: (instance) => (instance.activation ? "Y" : ""),
  };

  return [SELECTION, NAME, ENVIRONMENT, ADDRESS, EXTERNAL_LINK, LICENSE].filter(
    (column) => !column.hide
  );
});

const handleDataSourceToggle = (e: MouseEvent, instance: ComposedInstance) => {
  e.stopPropagation();
  if (state.dataSourceToggle.has(instance.name)) {
    state.dataSourceToggle.delete(instance.name);
  } else {
    state.dataSourceToggle.add(instance.name);
  }
};

const rowProps = (instance: ComposedInstance) => {
  return {
    style: "cursor: pointer;",
    onClick: (e: MouseEvent) => {
      if (props.onClick) {
        props.onClick(instance, e);
        return;
      }
      const url = `/${instance.name}`;
      if (e.ctrlKey || e.metaKey) {
        window.open(url, "_blank");
      } else {
        router.push(url);
      }
    },
  };
};
</script>
