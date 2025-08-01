import { create as createProto } from "@bufbuild/protobuf";
import { t } from "@/plugins/i18n";
import {
  Activity_Type,
  WebhookSchema,
  Webhook_Type,
} from "../proto-es/v1/project_service_pb";

export const emptyProjectWebhook = () => {
  return createProto(WebhookSchema, {
    type: Webhook_Type.SLACK,
    notificationTypes: [Activity_Type.ISSUE_STATUS_UPDATE],
  });
};

type ProjectWebhookV1TypeItem = {
  type: Webhook_Type;
  name: string;
  urlPrefix: string;
  urlPlaceholder: string;
  docUrl: string;
  supportDirectMessage: boolean;
};

export const projectWebhookV1TypeItemList = (): ProjectWebhookV1TypeItem[] => {
  return [
    {
      type: Webhook_Type.SLACK,
      name: t("common.slack"),
      urlPrefix: "https://hooks.slack.com/",
      urlPlaceholder: "https://hooks.slack.com/services/...",
      docUrl: "https://api.slack.com/messaging/webhooks",
      supportDirectMessage: true,
    },
    {
      type: Webhook_Type.DISCORD,
      name: t("common.discord"),
      urlPrefix: "https://discord.com/api/webhooks",
      urlPlaceholder: "https://discord.com/api/webhooks/...",
      docUrl:
        "https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks",
      supportDirectMessage: false,
    },
    {
      type: Webhook_Type.TEAMS,
      name: t("common.teams"),
      urlPrefix: "",
      urlPlaceholder: "https://acme123.webhook.office.com/webhookb2/...",
      docUrl:
        "https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook",
      supportDirectMessage: false,
    },
    {
      type: Webhook_Type.DINGTALK,
      name: t("common.dingtalk"),
      urlPrefix: "https://oapi.dingtalk.com",
      urlPlaceholder: "https://oapi.dingtalk.com/robot/...",
      docUrl:
        "https://developers.dingtalk.com/document/robots/custom-robot-access",
      supportDirectMessage: true,
    },
    {
      type: Webhook_Type.FEISHU,
      name: t("common.feishu"),
      urlPrefix: "https://open.feishu.cn",
      urlPlaceholder: "https://open.feishu.cn/open-apis/bot/v2/hook/...",
      docUrl:
        "https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot",
      supportDirectMessage: true,
    },
    {
      type: Webhook_Type.LARK,
      name: t("common.lark"),
      urlPrefix: "https://open.larksuite.com",
      urlPlaceholder: "https://open.larksuite.com/open-apis/bot/v2/hook/...",
      docUrl:
        "https://open.larksuite.com/document/client-docs/bot-v3/add-custom-bot",
      supportDirectMessage: true,
    },
    {
      type: Webhook_Type.WECOM,
      name: t("common.wecom"),
      urlPrefix: "https://qyapi.weixin.qq.com",
      urlPlaceholder: "https://qyapi.weixin.qq.com/cgi-bin/webhook/...",
      docUrl: "https://open.work.weixin.qq.com/help2/pc/14931",
      supportDirectMessage: true,
    },
  ];
};

type ProjectWebhookV1ActivityItem = {
  title: string;
  label: string;
  activity: Activity_Type;
  supportDirectMessage: boolean;
};

export const projectWebhookV1ActivityItemList =
  (): ProjectWebhookV1ActivityItem[] => {
    return [
      {
        title: t("project.webhook.activity-item.issue-creation.title"),
        label: t("project.webhook.activity-item.issue-creation.label"),
        activity: Activity_Type.ISSUE_CREATE,
        supportDirectMessage: false,
      },
      {
        title: t("project.webhook.activity-item.issue-status-change.title"),
        label: t("project.webhook.activity-item.issue-status-change.label"),
        activity: Activity_Type.ISSUE_STATUS_UPDATE,
        supportDirectMessage: false,
      },
      {
        title: t(
          "project.webhook.activity-item.issue-stage-status-change.title"
        ),
        label: t(
          "project.webhook.activity-item.issue-stage-status-change.label"
        ),
        activity: Activity_Type.ISSUE_PIPELINE_STAGE_STATUS_UPDATE,
        supportDirectMessage: false,
      },
      {
        title: t(
          "project.webhook.activity-item.issue-task-status-change.title"
        ),
        label: t(
          "project.webhook.activity-item.issue-task-status-change.label"
        ),
        activity: Activity_Type.ISSUE_PIPELINE_TASK_RUN_STATUS_UPDATE,
        supportDirectMessage: false,
      },
      {
        title: t("project.webhook.activity-item.issue-info-change.title"),
        label: t("project.webhook.activity-item.issue-info-change.label"),
        activity: Activity_Type.ISSUE_FIELD_UPDATE,
        supportDirectMessage: false,
      },
      {
        title: t("project.webhook.activity-item.issue-comment-creation.title"),
        label: t("project.webhook.activity-item.issue-comment-creation.label"),
        activity: Activity_Type.ISSUE_COMMENT_CREATE,
        supportDirectMessage: false,
      },
      {
        title: t("project.webhook.activity-item.issue-approval-notify.title"),
        label: t("project.webhook.activity-item.issue-approval-notify.label"),
        activity: Activity_Type.ISSUE_APPROVAL_NOTIFY,
        supportDirectMessage: true,
      },
      {
        title: t("project.webhook.activity-item.notify-issue-approved.title"),
        label: t("project.webhook.activity-item.notify-issue-approved.label"),
        activity: Activity_Type.NOTIFY_ISSUE_APPROVED,
        supportDirectMessage: true,
      },
      {
        title: t("project.webhook.activity-item.notify-pipeline-rollout.title"),
        label: t("project.webhook.activity-item.notify-pipeline-rollout.label"),
        activity: Activity_Type.NOTIFY_PIPELINE_ROLLOUT,
        supportDirectMessage: true,
      },
    ];
  };
