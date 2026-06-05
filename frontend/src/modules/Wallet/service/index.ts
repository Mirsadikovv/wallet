import { api } from "@/plugins/axios.plugin";
import { Try } from "@/common/try";

export type WalletType = {
  id: number;
  address: string;
  wallet_type: string;
  network: string;
  is_active: boolean;
  created_at: string;
};

export type WalletInfoType = {
  id: number;
  address: string;
  balance: string;
  wallet_type: string;
  network: string;
  seqno: number;
  is_active: boolean;
  created_at: string;
};

export type WalletCreateType = {
  user_id: number;
  wallet_type: string;
  network: string;
};

export type WalletListResponse = {
  wallets: WalletType[];
  total: number;
};

export type TransactionType = {
  hash: string;
  lt: number;
  timestamp: number;
  type: string;
  amount: string;
  fee: string;
  from: string;
  to: string;
  comment?: string;
  success: boolean;
};

export type TransactionsResponse = {
  wallet_id: number;
  address: string;
  transactions: TransactionType[];
  total: number;
};

export type SendCoinsType = {
  recipient: string;
  amount: string;
  comment?: string;
};

export type SendCoinsResponse = {
  hash: string;
  lt: number;
  address: string;
  amount: string;
  fee: string;
  recipient: string;
  comment?: string;
};

export type BalanceType = {
  address: string;
  balance: string;
};

const errHandler = async (err: unknown) => {
  const { ErrorNotify } = await import("@/common/Notify");
  ErrorNotify(
    (err as { response?: { data?: { message?: string } }; message: string })?.response?.data
      ?.message || (err as Error).message,
  );
};

class WalletService {
  @Try({ onError: errHandler })
  async list(user_id: number) {
    const { data } = await api.get<WalletListResponse>(`/v1/wallet/list?user_id=${user_id}`);
    return data;
  }

  @Try({ onError: errHandler })
  async create(dto: WalletCreateType) {
    const { data } = await api.post<WalletType>("/v1/wallet", dto);
    return data;
  }

  @Try({ onError: errHandler })
  async getInfo(id: number) {
    const { data } = await api.get<WalletInfoType>(`/v1/wallet/${id}`);
    return data;
  }

  @Try({ onError: errHandler })
  async getBalance(id: number) {
    const { data } = await api.get<BalanceType>(`/v1/wallet/${id}/balance`);
    return data;
  }

  @Try({ onError: errHandler })
  async getTransactions(id: number, limit = 10) {
    const { data } = await api.get<TransactionsResponse>(
      `/v1/wallet/${id}/transactions?limit=${limit}`,
    );
    return data;
  }

  @Try({ onError: errHandler })
  async send(id: number, dto: SendCoinsType) {
    const { data } = await api.post<SendCoinsResponse>(`/v1/wallet/${id}/send`, dto);
    return data;
  }

  @Try({ onError: errHandler })
  async remove(id: number) {
    const { data } = await api.delete<{ success: boolean; message: string }>(`/v1/wallet/${id}`);
    return data;
  }
}

const walletService = new WalletService();
export { walletService as WalletService };
