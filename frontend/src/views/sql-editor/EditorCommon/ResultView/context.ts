import { type Table } from "@tanstack/vue-table";
import { type Ref, type InjectionKey, provide, inject } from "vue";
import type { QueryRow } from "@/types/proto-es/v1/sql_service_pb";

export type SQLResultViewContext = {
  dark: Ref<boolean>;
  disallowCopyingData: Ref<boolean>;
  keyword: Ref<string>;
  detail: Ref<
    | {
        set: number; // The index of selected result set.
        row: number; // The row index of selected record.
        col: number; // The column index of selected cell.
        table: Table<QueryRow>;
      }
    | undefined
  >;
};

export const KEY = Symbol(
  "bb.sql-editor.result-view"
) as InjectionKey<SQLResultViewContext>;

export const provideSQLResultViewContext = (context: SQLResultViewContext) => {
  provide(KEY, context);
};

export const useSQLResultViewContext = () => {
  return inject(KEY)!;
};
