<template>
  <div class="flex items-start gap-3 px-3 py-2 border rounded-lg bg-gray-50">
    <component
      :is="statusIcon"
      class="w-5 h-5 flex-shrink-0"
      :class="statusColor"
    />

    <div class="flex-1 min-w-0 space-y-1">
      <div class="text-sm font-medium text-main">
        {{ displayTitle }}
      </div>
      <div v-if="content" class="text-sm text-control">
        {{ content }}
      </div>
      <div v-if="position" class="text-sm mt-1 text-control-light">
        Line {{ position.line + 1 }}, Column {{ position.column + 1 }}
      </div>
      <div
        v-if="affectedRows !== undefined"
        class="text-sm mt-1 flex items-center gap-1"
      >
        <NTag size="small" round>
          {{ $t("task.check-type.affected-rows.self") }}
        </NTag>
        <span>{{ affectedRows }}</span>
        <span class="text-control opacity-80"
          >({{ $t("task.check-type.affected-rows.description") }})</span
        >
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CheckCircleIcon, AlertCircleIcon, XCircleIcon } from "lucide-vue-next";
import { NTag } from "naive-ui";
import { computed } from "vue";
import { getRuleLocalization, ruleTemplateMapV2 } from "@/types";

export interface CheckResultPosition {
  line: number;
  column: number;
}

const props = defineProps<{
  status: "SUCCESS" | "WARNING" | "ERROR";
  title: string;
  content?: string;
  position?: CheckResultPosition;
  affectedRows?: number;
  code?: number;
  reportType?: "sqlReviewReport" | "sqlSummaryReport" | string;
}>();

const statusIcon = computed(() => {
  switch (props.status) {
    case "ERROR":
      return XCircleIcon;
    case "WARNING":
      return AlertCircleIcon;
    case "SUCCESS":
      return CheckCircleIcon;
    default:
      return CheckCircleIcon;
  }
});

const statusColor = computed(() => {
  switch (props.status) {
    case "ERROR":
      return "text-error";
    case "WARNING":
      return "text-warning";
    case "SUCCESS":
      return "text-success";
    default:
      return "text-control";
  }
});

const getRuleTemplateByType = (type: string) => {
  for (const mapByType of ruleTemplateMapV2.values()) {
    if (mapByType.has(type)) {
      return mapByType.get(type);
    }
  }
  return;
};

const messageWithCode = (message: string, code: number | undefined): string => {
  if (code !== undefined && code !== 0) {
    return `${message} #${code}`;
  }
  return message;
};

const displayTitle = computed(() => {
  let title = props.title;

  // Skip localization for certain titles
  if (title === "OK" || title === "Syntax error") {
    return messageWithCode(title, props.code);
  }

  // Only apply SQL review localization if this is a SQL review report
  if (props.reportType === "sqlReviewReport") {
    // Convert dots to hyphens in the rule key to match the expected format
    const normalizedKey = title.replace(/\./g, "-");
    // Try to get rule template and localization
    const rule = getRuleTemplateByType(props.title);
    if (rule) {
      const ruleLocalization = getRuleLocalization(rule.type, rule.engine);
      if (ruleLocalization.title) {
        title = ruleLocalization.title;
      }
    } else {
      // Fallback to direct localization lookup
      const localization = getRuleLocalization(normalizedKey);
      if (localization.title) {
        title = localization.title;
      }
    }
  }

  // Add error code if present
  return messageWithCode(title, props.code);
});
</script>
