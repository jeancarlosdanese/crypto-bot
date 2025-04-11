// File: types/paginator.ts

export type Paginator<T> = {
  total_records: number;
  total_pages: number;
  current_page: number;
  per_page: number;
  data: T[];
};
