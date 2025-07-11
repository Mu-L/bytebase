<template>
  <div class="relative">
    <div v-if="icon === 'system'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-control-bg rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <img class="mt-1" src="@/assets/logo-icon.svg" alt="Bytebase" />
      </div>
    </div>
    <div v-else-if="icon === 'avatar'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-white rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <PrincipalAvatar
          :user="user"
          override-class="w-7 h-7 font-medium"
          override-text-size="0.8rem"
        />
      </div>
    </div>
    <div v-else-if="icon === 'create'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-control-bg rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <PlusIcon class="w-5 h-5 text-control" />
      </div>
    </div>
    <div v-else-if="icon === 'update'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-control-bg rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <PencilIcon class="w-4 h-4 text-control" />
      </div>
    </div>
    <div
      v-else-if="icon === 'run' || icon === 'rollout'"
      class="relative pl-0.5"
    >
      <div
        class="w-7 h-7 bg-control-bg rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <PlayCircleIcon class="w-5 h-5 text-control" />
      </div>
    </div>
    <div v-else-if="icon === 'approve-review'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-success rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <ThumbsUpIcon class="w-4 h-4 text-white" />
      </div>
    </div>
    <div v-else-if="icon === 'reject-review'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-warning rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <PencilIcon class="w-4 h-4 text-white" />
      </div>
    </div>
    <div v-else-if="icon === 're-request-review'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-control-bg rounded-full ring-4 ring-white flex items-center justify-center icon-re-request-review"
      >
        <PlayIcon class="w-4 h-4 text-control ml-px" />
      </div>
    </div>
    <div v-else-if="icon === 'cancel'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-control-bg rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <MinusIcon class="w-5 h-5 text-control" />
      </div>
    </div>
    <div v-else-if="icon === 'fail'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-white rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <CircleAlertIcon class="w-5 h-5 text-error" />
      </div>
    </div>
    <div v-else-if="icon === 'complete'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-white rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <CheckCircle2Icon class="w-6 h-6 text-success" />
      </div>
    </div>
    <div v-else-if="icon === 'skip'" class="relative pl-1">
      <div
        class="w-6 h-6 bg-gray-200 rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <SkipIcon class="w-5 h-5 text-gray-500" />
      </div>
    </div>
    <div v-else-if="icon === 'commit'" class="relative pl-0.5">
      <div
        class="w-7 h-7 bg-control-bg rounded-full ring-4 ring-white flex items-center justify-center"
      >
        <CodeIcon class="w-5 h-5 text-control" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computedAsync } from "@vueuse/core";
import {
  ThumbsUpIcon,
  PlayIcon,
  CheckCircle2Icon,
  PlusIcon,
  PencilIcon,
  PlayCircleIcon,
  MinusIcon,
  CodeIcon,
  CircleAlertIcon,
} from "lucide-vue-next";
import { computed } from "vue";
import { SkipIcon } from "@/components/Icon";
import PrincipalAvatar from "@/components/PrincipalAvatar.vue";
import {
  IssueCommentType,
  useUserStore,
  type ComposedIssueComment,
} from "@/store";
import { extractUserId } from "@/store";
import {
  IssueComment_Approval_Status,
  IssueComment_TaskUpdate_Status,
} from "@/types/proto-es/v1/issue_service_pb";

type ActionIconType =
  | "avatar"
  | "system"
  | "create"
  | "update"
  | "run"
  | "approve-review"
  | "reject-review"
  | "re-request-review"
  | "rollout"
  | "cancel"
  | "fail"
  | "complete"
  | "skip"
  | "commit";

const props = defineProps<{
  issueComment: ComposedIssueComment;
}>();

const userStore = useUserStore();

const user = computedAsync(() => {
  return userStore.getOrFetchUserByIdentifier(props.issueComment.creator);
});

const icon = computed((): ActionIconType => {
  const { issueComment } = props;
  if (
    issueComment.type === IssueCommentType.APPROVAL &&
    issueComment.event?.case === "approval"
  ) {
    const { status } = issueComment.event.value;
    switch (status) {
      case IssueComment_Approval_Status.APPROVED:
        return "approve-review";
      case IssueComment_Approval_Status.REJECTED:
        return "reject-review";
      case IssueComment_Approval_Status.PENDING:
        return "re-request-review";
    }
  } else if (issueComment.type === IssueCommentType.STAGE_END) {
    return "complete";
  } else if (
    issueComment.type === IssueCommentType.TASK_UPDATE &&
    issueComment.event?.case === "taskUpdate"
  ) {
    const { toStatus } = issueComment.event.value;
    let action: ActionIconType = "update";
    if (toStatus !== undefined) {
      switch (toStatus) {
        case IssueComment_TaskUpdate_Status.PENDING: {
          action = "rollout";
          break;
        }
        case IssueComment_TaskUpdate_Status.RUNNING: {
          action = "run";
          break;
        }
        case IssueComment_TaskUpdate_Status.DONE: {
          action = "complete";
          break;
        }
        case IssueComment_TaskUpdate_Status.FAILED: {
          action = "fail";
          break;
        }
        case IssueComment_TaskUpdate_Status.SKIPPED: {
          action = "skip";
          break;
        }
        case IssueComment_TaskUpdate_Status.CANCELED: {
          action = "cancel";
          break;
        }
      }
    }
    return action;
  } else if (
    issueComment.type === IssueCommentType.ISSUE_UPDATE &&
    issueComment.event?.case === "issueUpdate"
  ) {
    const { toTitle, toDescription, toLabels, fromLabels } =
      issueComment.event.value;
    if (
      toTitle !== undefined ||
      toDescription !== undefined ||
      toLabels.length !== 0 ||
      fromLabels.length !== 0
    ) {
      return "update";
    }
    // Otherwise, show avatar icon based on the creator.
  } else if (
    issueComment.type === IssueCommentType.TASK_PRIOR_BACKUP &&
    issueComment.event?.case === "taskPriorBackup"
  ) {
    const taskPriorBackup = issueComment.event.value;
    if (taskPriorBackup.error !== "") {
      return "fail";
    } else {
      return "complete";
    }
  }

  return extractUserId(issueComment.creator) == userStore.systemBotUser?.email
    ? "system"
    : "avatar";
});
</script>

<style scoped>
.icon-re-request-review :deep(path) {
  stroke-width: 3 !important;
}
</style>
