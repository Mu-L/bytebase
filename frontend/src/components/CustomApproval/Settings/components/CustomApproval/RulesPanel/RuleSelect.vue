<template>
  <div
    class="flex items-center"
    :class="attrs.class as VueClass"
    :style="attrs.style as VueStyle"
  >
    <label v-if="label" class="mr-2">{{ label }}</label>
    <SpinnerSelect
      :value="value"
      :on-update="onUpdate"
      :options="options"
      :placeholder="$t('custom-approval.approval-flow.select')"
      :consistent-menu-width="false"
      :disabled="disabled || !allowAdmin"
      :filterable="true"
      class="bb-rule-select"
      v-bind="selectAttrs"
    />
    <NButton
      v-if="link"
      quaternary
      type="info"
      class="!rounded !w-[var(--n-height)] !p-0 !ml-1"
      :disabled="!selectedRule || !allowAdmin"
      @click="toApprovalFlow"
    >
      <heroicons:pencil-square class="w-5 h-5" />
    </NButton>
  </div>
</template>

<script lang="ts" setup>
import { omit } from "lodash-es";
import type { SelectOption } from "naive-ui";
import { NButton, type SelectProps } from "naive-ui";
import { computed, useAttrs } from "vue";
import { useI18n } from "vue-i18n";
import { SpinnerSelect } from "@/components/v2/Form";
import { useWorkspaceApprovalSettingStore } from "@/store";
import type { VueClass, VueStyle } from "@/utils";
import { useCustomApprovalContext } from "../context";

export interface ApprovalTemplateSelectorProps
  extends /* @vue-ignore */ SelectProps {
  label?: string;
  link?: boolean;
  disabled?: boolean;
  value?: string;
  onUpdate: (value: string | undefined) => Promise<any>;
  selectClass?: VueClass;
  selectStyle?: VueStyle;
}
const props = defineProps<ApprovalTemplateSelectorProps>();

defineEmits<{
  (event: "update:value", value: string | undefined): void;
}>();

const { t } = useI18n();
const store = useWorkspaceApprovalSettingStore();
const context = useCustomApprovalContext();
const { allowAdmin, hasFeature, showFeatureModal } = context;

const attrs = useAttrs();
const selectAttrs = computed(() => ({
  ...omit(attrs, "class", "style"),
  class: props.selectClass,
  style: props.selectStyle,
}));

const options = computed(() => {
  const ruleOptions = store.config.rules.map<SelectOption>((rule) => ({
    label: rule.template.title,
    value: rule.uid,
  }));
  return [
    { value: "", label: t("custom-approval.approval-flow.skip") },
    ...ruleOptions,
  ];
});

const selectedRule = computed(() => {
  return store.config.rules.find((rule) => rule.uid === props.value);
});

const toApprovalFlow = () => {
  const rule = selectedRule.value;
  if (!rule) {
    return;
  }
  if (!hasFeature.value) {
    showFeatureModal.value = true;
    return;
  }
  context.dialog.value = {
    mode: "EDIT",
    rule,
  };
};
</script>
