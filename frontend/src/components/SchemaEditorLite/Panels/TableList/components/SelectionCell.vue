<template>
  <NCheckbox
    :checked="state.checked"
    :indeterminate="state.indeterminate"
    @update:checked="update"
  />
</template>

<script setup lang="ts">
import { NCheckbox } from "naive-ui";
import { computed } from "vue";
import { useSchemaEditorContext } from "@/components/SchemaEditorLite/context";
import type { ComposedDatabase } from "@/types";
import type {
  DatabaseMetadata,
  SchemaMetadata,
  TableMetadata,
} from "@/types/proto-es/v1/database_service_pb";

const props = defineProps<{
  db: ComposedDatabase;
  metadata: {
    database: DatabaseMetadata;
    schema: SchemaMetadata;
    table: TableMetadata;
  };
}>();
const { getTableSelectionState, updateTableSelection } =
  useSchemaEditorContext();

const state = computed(() => {
  return getTableSelectionState(props.db, props.metadata);
});
const update = (on: boolean) => {
  updateTableSelection(props.db, props.metadata, on);
};
</script>
