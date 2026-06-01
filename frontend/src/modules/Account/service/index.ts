import { api } from "@/plugins/axios.plugin";
import { Try } from "@/common/try";

export type AccountType = {
  id: number;
  user_id: number;
  name: string;
  currency: string;
  balance: string;
  is_active: boolean;
  created_at: string;
};

export type AccountCreateType = {
  user_id: number;
  name: string;
  currency: string;
};

export type AccountUpdateType = {
  name: string;
};

export type AccountListResponse = {
  accounts: AccountType[];
  total: number;
};

const errHandler = async (err: unknown) => {
  const { ErrorNotify } = await import("@/common/Notify");
  ErrorNotify(
    (err as { response?: { data?: { message?: string } }; message: string })?.response?.data
      ?.message || (err as Error).message,
  );
};

class AccountService {
  @Try({ onError: errHandler })
  async list(user_id: number) {
    const { data } = await api.get<AccountListResponse>(`/v1/account/list?user_id=${user_id}`);
    return data;
  }

  @Try({ onError: errHandler })
  async create(dto: AccountCreateType) {
    const { data } = await api.post<AccountType>("/v1/account", dto);
    return data;
  }

  @Try({ onError: errHandler })
  async getByID(id: number) {
    const { data } = await api.get<AccountType>(`/v1/account/${id}`);
    return data;
  }

  @Try({ onError: errHandler })
  async update(id: number, dto: AccountUpdateType) {
    const { data } = await api.put<AccountType>(`/v1/account/${id}`, dto);
    return data;
  }

  @Try({ onError: errHandler })
  async deleteOrRestore(id: number) {
    const { data } = await api.delete<{ success: boolean; message: string }>(
      `/v1/account/${id}`,
    );
    return data;
  }
}

const accountService = new AccountService();
export { accountService as AccountService };
