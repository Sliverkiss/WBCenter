import { FormEvent, useCallback, useEffect, useState } from "react";

type Page =
  | "overview"
  | "accounts"
  | "credits"
  | "models"
  | "automation"
  | "scheduler"
  | "activities"
  | "oauth";
type Any = Record<string, unknown>;
type ActivityRun = {
  id: string;
  uid: string;
  nickname: string;
  status: "running" | "completed" | "failed";
  total: number;
  completed: number;
  succeeded: number;
  already: number;
  skipped: number;
  failed: number;
  credit: number;
  energy: number;
  duration_seconds: number;
  summary: string;
  logs: string[];
};

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    credentials: "same-origin",
    headers: init?.body ? { "Content-Type": "application/json" } : undefined,
    ...init,
  });
  const body = (await response.json().catch(() => ({}))) as Any;
  if (!response.ok) throw new Error(String(body.error || "请求失败"));
  return body as T;
}
const get = <T,>(url: string) => request<T>(url);
const post = <T,>(url: string, body: unknown = {}) =>
  request<T>(url, { method: "POST", body: JSON.stringify(body) });

export default function App() {
  const [authed, setAuthed] = useState<boolean | null>(null);
  const [readOnly, setReadOnly] = useState(false);
  const [page, setPage] = useState<Page>("overview");
  const [notice, setNotice] = useState("");

  const check = useCallback(async () => {
    try {
      const s = await get<{ authenticated: boolean; read_only?: boolean }>(
        "/api/session",
      );
      setAuthed(s.authenticated);
      setReadOnly(Boolean(s.read_only));
    } catch {
      setAuthed(false);
    }
  }, []);
  useEffect(() => {
    void check();
  }, [check]);
  useEffect(() => {
    if (!notice) return;
    const timer = window.setTimeout(() => setNotice(""), 6000);
    return () => window.clearTimeout(timer);
  }, [notice]);
  if (authed === null) return <div className="center">正在连接本地控制台…</div>;
  if (!authed) return <Login onDone={check} />;

  const nav: [Page, string, string][] = [
    ["overview", "概览", "⌁"],
    ["accounts", "账号管理", "◉"],
    ["credits", "积分管理", "＋"],
    ["models", "模型管理", "◇"],
    ["automation", "自动化管理", "◷"],
    ["scheduler", "账号云端定时任务", "☁"],
    ["activities", "活动管理", "✦"],
    ["oauth", "添加账号", "→"],
  ];
  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="wordmark">
          <b>
            WORK
            <br />
            BUDDY<span>›</span>
          </b>
          <small>LOCAL CONTROL ROOM</small>
        </div>
        <nav>
          {nav.map(([id, label, icon]) => (
            <button
              key={id}
              className={page === id ? "nav active" : "nav"}
              onClick={() => setPage(id)}
            >
              <i>{icon}</i>
              {label}
            </button>
          ))}
        </nav>
        <div className="side-bottom">
          <span>第三方本地控制台</span>
          <button
            className="link"
            onClick={() =>
              void post("/api/logout").then(() => setAuthed(false))
            }
          >
            退出登录
          </button>
        </div>
      </aside>
      <main className="main">
        {readOnly && (
          <div className="notice neutral">
            当前为只读模式，账号与配置写操作已禁用。
          </div>
        )}
        {notice && (
          <div className="notice">
            {notice}
            <button onClick={() => setNotice("")}>×</button>
          </div>
        )}
        {page === "overview" && <Overview onPage={setPage} />}
        {page === "accounts" && (
          <Accounts onNotice={setNotice} readOnly={readOnly} />
        )}
        {page === "credits" && <Credits />}
        {page === "models" && <Models />}
        {page === "automation" && (
          <Automation onNotice={setNotice} readOnly={readOnly} />
        )}
        {page === "scheduler" && <Scheduler readOnly={readOnly} />}
        {page === "activities" && (
          <Activities onNotice={setNotice} readOnly={readOnly} />
        )}
        {page === "oauth" && (
          <OAuth
            onNotice={setNotice}
            onDone={() => setPage("accounts")}
            readOnly={readOnly}
          />
        )}
      </main>
    </div>
  );
}

function Login({ onDone }: { onDone: () => Promise<void> }) {
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const submit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await post("/api/login", { username, password });
      await onDone();
    } catch (err) {
      setError(err instanceof Error ? err.message : "登录失败");
    }
  };
  return (
    <div className="login">
      <div className="login-card">
        <div className="eyebrow">LOCAL CONTROL ROOM</div>
        <h1>账号工作台</h1>
        <p>只读取账号凭据并直接连接 WorkBuddy 上游。</p>
        <form onSubmit={submit}>
          <label>
            用户名
            <input
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoComplete="username"
            />
          </label>
          <label>
            密码
            <input
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              autoFocus
              autoComplete="current-password"
            />
          </label>
          {error && <div className="error">{error}</div>}
          <button className="black" type="submit">
            进入控制台
          </button>
        </form>
      </div>
    </div>
  );
}

function Head({
  eyebrow,
  title,
  text,
  action,
}: {
  eyebrow: string;
  title: string;
  text: string;
  action?: React.ReactNode;
}) {
  return (
    <header className="head">
      <div>
        <div className="eyebrow">{eyebrow}</div>
        <h1>{title}</h1>
        <p>{text}</p>
      </div>
      {action}
    </header>
  );
}
function Loading() {
  return <div className="loading">正在直接查询 WorkBuddy 上游…</div>;
}
function ErrorView({ text }: { text: string }) {
  return <div className="error card">{text}</div>;
}
function fmt(v: unknown) {
  const n = Number(v);
  return Number.isFinite(n)
    ? n.toLocaleString("zh-CN", { maximumFractionDigits: 2 })
    : "—";
}
function displayTime(v: unknown) {
  const raw = String(v || "");
  return raw && !raw.startsWith("0001-")
    ? new Date(raw).toLocaleString()
    : "未执行";
}

function Overview({ onPage }: { onPage: (p: Page) => void }) {
  const [data, setData] = useState<Any | null>(null);
  const [error, setError] = useState("");
  const load = useCallback(
    () =>
      void get<Any>("/api/overview")
        .then(setData)
        .catch((e) => setError(e.message)),
    [],
  );
  useEffect(load, [load]);
  if (error) return <ErrorView text={error} />;
  if (!data) return <Loading />;
  const cards = [
    ["健康账号", data.healthy],
    ["账号总数", data.accounts],
    ["已过期", data.expired],
    ["已启用自动化", data.automations],
  ];
  return (
    <>
      <Head
        eyebrow="LOCAL CONTROL ROOM"
        title="账号工作台"
        text="账号池状态、今日积分与已启用自动化"
        action={
          <button className="black" onClick={load}>
            刷新状态
          </button>
        }
      />
      <div className="stats">
        {cards.map(([label, value]) => (
          <div className="stat" key={String(label)}>
            <span>{String(label)}</span>
            <b>{String(value)}</b>
          </div>
        ))}
      </div>
      <div className="grid two">
        <section className="card">
          <h2>快速入口</h2>
          <div className="quick">
            <button onClick={() => onPage("credits")}>
              查看今日积分汇总 →
            </button>
            <button onClick={() => onPage("automation")}>
              管理本地自动化 →
            </button>
            <button onClick={() => onPage("scheduler")}>
              查看账号云端定时任务 →
            </button>
            <button onClick={() => onPage("activities")}>
              探测新活动任务 →
            </button>
          </div>
        </section>
        <section className="card">
          <h2>运行边界</h2>
          <p className="muted">
            此面板不会读取或改写 workbuddy2api 的配置，也不会调用 Docker
            或重启网关。凭据仅保存在共享的 auths 目录。
          </p>
        </section>
      </div>
    </>
  );
}

function Accounts({
  onNotice,
  readOnly,
}: {
  onNotice: (s: string) => void;
  readOnly: boolean;
}) {
  const [items, setItems] = useState<Any[] | null>(null);
  const [warnings, setWarnings] = useState<string[]>([]);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState("");
  const load = useCallback(() => {
    setError("");
    void get<{ accounts: Any[]; warnings?: string[] }>("/api/accounts")
      .then((d) => {
        setItems(d.accounts);
        setWarnings(d.warnings || []);
      })
      .catch((e) => setError(e.message));
  }, []);
  useEffect(load, [load]);
  const action = async (uid: string, a: string, name: string) => {
    const labels: Record<string, string> = {
      refresh: "刷新凭据",
      checkin: "签到",
      travel: "旅行巡检",
    };
    if (
      a !== "refresh" &&
      !window.confirm(`确认立即为「${name}」执行${labels[a]}？`)
    )
      return;
    const key = `${uid}:${a}`;
    setBusy(key);
    try {
      const r = await post<{ message: string }>(
        `/api/accounts/${encodeURIComponent(uid)}/actions/${a}`,
      );
      onNotice(`${name}：${r.message}`);
      load();
    } catch (e) {
      onNotice(e instanceof Error ? e.message : "执行失败");
    } finally {
      setBusy("");
    }
  };
  if (error) return <ErrorView text={error} />;
  if (!items) return <Loading />;
  return (
    <>
      <Head
        eyebrow="ACCOUNTS"
        title="账号管理"
        text="OAuth 添加、凭据刷新与单账号操作"
        action={
          <button className="black" onClick={load}>
            刷新
          </button>
        }
      />
      {warnings.map((w, i) => (
        <div className="notice error" key={i}>
          {w}
        </div>
      ))}
      <section className="card table">
        <table>
          <thead>
            <tr>
              <th>账号</th>
              <th>区域</th>
              <th>Token</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {items.map((a) => {
              const uid = String(a.uid);
              const name = String(a.nickname || uid);
              return (
                <tr key={uid}>
                  <td>
                    <b>{name}</b>
                    <small>{uid}</small>
                  </td>
                  <td>{String(a.domain || "cn")}</td>
                  <td>
                    <Badge ok={!a.expired}>
                      {a.expired
                        ? "已过期"
                        : a.needs_refresh
                          ? "即将过期"
                          : "正常"}
                    </Badge>
                  </td>
                  <td>
                    <button
                      disabled={readOnly || Boolean(busy)}
                      onClick={() => void action(uid, "refresh", name)}
                    >
                      {busy === `${uid}:refresh` ? "执行中…" : "刷新凭据"}
                    </button>
                    <button
                      disabled={readOnly || Boolean(busy)}
                      onClick={() => void action(uid, "checkin", name)}
                    >
                      {busy === `${uid}:checkin` ? "执行中…" : "签到"}
                    </button>
                    <button
                      disabled={readOnly || Boolean(busy)}
                      onClick={() => void action(uid, "travel", name)}
                    >
                      {busy === `${uid}:travel` ? "执行中…" : "旅行巡检"}
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        {items.length === 0 && (
          <Empty text="尚无账号，请先通过 OAuth 添加账号。" />
        )}
      </section>
    </>
  );
}

function Credits() {
  const [items, setItems] = useState<Any[] | null>(null);
  const [error, setError] = useState("");
  const [detail, setDetail] = useState<Any | null>(null);
  const [detailBusy, setDetailBusy] = useState("");
  const load = useCallback(() => {
    setError("");
    void get<{ items: Any[] }>("/api/credits")
      .then((d) => setItems(d.items))
      .catch((e) => setError(e.message));
  }, []);
  useEffect(load, [load]);
  const inspect = async (x: Any) => {
    const uid = String(x.uid);
    setDetailBusy(uid);
    try {
      const d = await get<{ items: Any[] }>(
        `/api/credits/${encodeURIComponent(uid)}`,
      );
      setDetail(d.items[0] || null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "明细查询失败");
    } finally {
      setDetailBusy("");
    }
  };
  if (error) return <ErrorView text={error} />;
  if (!items) return <Loading />;
  return (
    <>
      <Head
        eyebrow="CREDITS"
        title="积分管理"
        text="默认展示今日额度发放、已用、剩余；不建立本地积分账本。"
        action={
          <button className="black" onClick={load}>
            查询上游
          </button>
        }
      />
      <section className="card table">
        <table>
          <thead>
            <tr>
              <th>账号</th>
              <th>当前积分</th>
              <th>今日发放额度</th>
              <th>今日已用</th>
              <th>今日剩余</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {items.map((x) => (
              <tr key={String(x.uid)}>
                <td>{String(x.nickname || x.uid)}</td>
                <td>
                  <b>{fmt(x.current)}</b>
                </td>
                <td>{fmt(x.today_allocated)}</td>
                <td>{fmt(x.today_consumed)}</td>
                <td>{fmt(x.today_remaining)}</td>
                <td>
                  <button
                    disabled={Boolean(detailBusy)}
                    onClick={() => void inspect(x)}
                  >
                    {detailBusy === String(x.uid) ? "查询中…" : "查看明细"}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
      <p className="muted">
        “今日发放/已用/剩余”来自上游免费套餐的当日切片，并非虚构的积分交易流水。
      </p>
      {detail && (
        <div className="drawer">
          <button
            className="close"
            aria-label="关闭积分明细"
            onClick={() => setDetail(null)}
          >
            ×
          </button>
          <h2>{String(detail.nickname)} 的今日额度明细</h2>
          <dl>
            <dt>当前积分</dt>
            <dd>{fmt(detail.current)}</dd>
            <dt>套餐数</dt>
            <dd>{fmt(detail.packages)}</dd>
            <dt>查询时间</dt>
            <dd>{new Date(String(detail.fetched_at)).toLocaleString()}</dd>
          </dl>
          {Boolean(detail.error) && (
            <div className="error">{String(detail.error)}</div>
          )}
          {Array.isArray(detail.package_details) &&
          detail.package_details.length > 0 ? (
            <div className="detail-packages">
              <h3>今日套餐切片</h3>
              <table>
                <thead>
                  <tr>
                    <th>套餐</th>
                    <th>发放</th>
                    <th>已用</th>
                    <th>剩余</th>
                  </tr>
                </thead>
                <tbody>
                  {detail.package_details.map((p: Any, i: number) => (
                    <tr key={`${String(p.code)}-${i}`}>
                      <td>
                        {String(p.name || p.code || "未命名套餐")}
                        <small>{String(p.code || "")}</small>
                      </td>
                      <td>{fmt(p.total)}</td>
                      <td>{fmt(p.used)}</td>
                      <td>{fmt(p.remaining)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            !detail.error && <Empty text="上游未返回今日套餐切片。" />
          )}
        </div>
      )}
    </>
  );
}

function Models() {
  const [items, setItems] = useState<Any[] | null>(null);
  const [error, setError] = useState("");
  const [query, setQuery] = useState("");
  const [account, setAccount] = useState("");
  const load = useCallback(() => {
    setError("");
    void get<{ items: Any[] }>("/api/models")
      .then((d) => setItems(d.items))
      .catch((e) => setError(e.message));
  }, []);
  useEffect(load, [load]);
  if (error) return <ErrorView text={error} />;
  if (!items) return <Loading />;
  const accounts = Array.from(
    new Map(
      items.map((m) => [String(m.uid), String(m.nickname || m.uid)]),
    ).entries(),
  );
  const q = query.trim().toLowerCase();
  const shown = items.filter(
    (m) =>
      (!account || String(m.uid) === account) &&
      (!q || `${String(m.id)} ${String(m.name)}`.toLowerCase().includes(q)),
  );
  return (
    <>
      <Head
        eyebrow="MODELS"
        title="模型管理"
        text="按账号直连上游查询可用模型，不依赖网关模型列表。"
        action={
          <button className="black" onClick={load}>
            刷新模型
          </button>
        }
      />
      <section className="card filters">
        <label>
          账号
          <select value={account} onChange={(e) => setAccount(e.target.value)}>
            <option value="">全部账号</option>
            {accounts.map(([uid, name]) => (
              <option key={uid} value={uid}>
                {name}
              </option>
            ))}
          </select>
        </label>
        <label>
          搜索模型
          <input
            type="search"
            value={query}
            placeholder="模型 ID 或名称"
            onChange={(e) => setQuery(e.target.value)}
          />
        </label>
        <span>{shown.length} 个结果</span>
      </section>
      <section className="card table">
        <table>
          <thead>
            <tr>
              <th>账号</th>
              <th>模型 ID</th>
              <th>名称</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            {shown.map((m, i) => (
              <tr key={`${String(m.uid)}-${i}`}>
                <td>{String(m.nickname || m.uid)}</td>
                <td className="mono">{String(m.id || "—")}</td>
                <td>{String(m.name || "—")}</td>
                <td>
                  {m.error ? (
                    <Badge ok={false}>{String(m.error)}</Badge>
                  ) : (
                    <Badge ok>来自上游</Badge>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {shown.length === 0 && <Empty text="没有符合筛选条件的模型。" />}
      </section>
    </>
  );
}

function Automation({
  onNotice,
  readOnly,
}: {
  onNotice: (s: string) => void;
  readOnly: boolean;
}) {
  const [items, setItems] = useState<Any[] | null>(null);
  const [runs, setRuns] = useState<Any[]>([]);
  const [error, setError] = useState("");
  const [intervals, setIntervals] = useState<Record<string, string>>({});
  const [busy, setBusy] = useState("");
  const load = useCallback(() => {
    setError("");
    void get<{ items: Any[]; runs: Any[] }>("/api/automations")
      .then((d) => {
        setItems(d.items);
        setRuns(d.runs);
        setIntervals(
          Object.fromEntries(
            d.items.map((a) => [String(a.id), String(a.every_minutes)]),
          ),
        );
      })
      .catch((e) => setError(e.message));
  }, []);
  useEffect(load, [load]);
  const every = (a: Any) => Number(intervals[String(a.id)] || a.every_minutes);
  const save = async (a: Any, enabled: boolean) => {
    const id = String(a.id);
    setBusy(id);
    try {
      await request(`/api/automations/${encodeURIComponent(id)}`, {
        method: "PUT",
        body: JSON.stringify({ enabled, every_minutes: every(a) }),
      });
      onNotice(`${String(a.name)}配置已保存`);
      load();
    } catch (e) {
      onNotice(e instanceof Error ? e.message : "保存失败");
    } finally {
      setBusy("");
    }
  };
  const run = async (a: Any) => {
    if (
      a.action !== "activity_probe" &&
      !window.confirm(`确认立即对全部账号执行「${String(a.name)}」？`)
    )
      return;
    const id = String(a.id);
    setBusy(id);
    try {
      const r = await post<{ message: string }>(
        `/api/automations/${encodeURIComponent(id)}/run`,
      );
      onNotice(`${String(a.name)}：${r.message}`);
      load();
    } catch (e) {
      onNotice(e instanceof Error ? e.message : "执行失败");
    } finally {
      setBusy("");
    }
  };
  if (error) return <ErrorView text={error} />;
  if (!items) return <Loading />;
  return (
    <>
      <Head
        eyebrow="AUTOMATION"
        title="自动化管理"
        text="这些是真实持久化的面板定时任务，可分别调整间隔、启停或立即运行。"
      />
      <div className="notice neutral">
        签到、旅行、保活和活动自动完成会影响真实账号。请避免其他服务重复执行同名任务。
      </div>
      <section className="card">
        <div className="automation-list">
          {items.map((a) => {
            const id = String(a.id);
            return (
              <div className="automation" key={id}>
                <div>
                  <h2>{String(a.name)}</h2>
                  <p>
                    下次：{a.enabled ? displayTime(a.next_run_at) : "已停用"} ·
                    上次：{displayTime(a.last_run_at)}
                  </p>
                  <small>{String(a.last_result || "等待执行")}</small>
                </div>
                <div className="actions">
                  <label className="interval">
                    每{" "}
                    <input
                      disabled={readOnly}
                      aria-label={`${String(a.name)}执行间隔`}
                      type="number"
                      min="5"
                      max="10080"
                      value={intervals[id] || String(a.every_minutes)}
                      onChange={(e) =>
                        setIntervals((v) => ({ ...v, [id]: e.target.value }))
                      }
                    />{" "}
                    分钟
                  </label>
                  <button
                    disabled={readOnly || Boolean(busy)}
                    onClick={() => void save(a, Boolean(a.enabled))}
                  >
                    保存配置
                  </button>
                  <button
                    disabled={readOnly || Boolean(busy)}
                    className={a.enabled ? "switch on" : "switch"}
                    aria-label={`切换${String(a.name)}`}
                    aria-pressed={Boolean(a.enabled)}
                    onClick={() => void save(a, !a.enabled)}
                  >
                    <span />
                  </button>
                  <button
                    disabled={readOnly || Boolean(busy)}
                    onClick={() => void run(a)}
                  >
                    {busy === id ? "执行中…" : "立即运行"}
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      </section>
      {runs.length > 0 && (
        <section className="card">
          <h2>最近运行</h2>
          {runs.slice(0, 8).map((r, i) => (
            <p className="run" key={i}>
              <Badge ok={Boolean(r.ok)}>{r.ok ? "成功" : "失败"}</Badge>{" "}
              {automationActionLabel(String(r.action))} · {String(r.message)} · {displayTime(r.at)}
            </p>
          ))}
        </section>
      )}
    </>
  );
}

function automationActionLabel(action: string) {
  const labels: Record<string, string> = {
    activity_probe: "新活动探测",
    activity_complete: "活动任务自动完成",
    checkin: "每日签到",
    travel: "旅行巡检",
    refresh: "账号保活",
  };
  return labels[action] || "自动化任务";
}

function Scheduler({ readOnly }: { readOnly: boolean }) {
  const [items, setItems] = useState<Any[] | null>(null);
  const [mock, setMock] = useState<Any[] | null>(null);
  const [accounts, setAccounts] = useState<Any[]>([]);
  const [note, setNote] = useState("");
  const [detail, setDetail] = useState<Any | null>(null);
  const [error, setError] = useState("");
  const [draft, setDraft] = useState({
    account_uid: "",
    name: "",
    cron: "0 9 * * *",
    prompt: "",
    enabled: true,
  });
  const [edits, setEdits] = useState<Record<string, Any>>({});
  const load = useCallback(
    () =>
      void Promise.all([
        get<{ items: Any[]; note: string }>("/api/scheduler-tasks"),
        get<{ items: Any[]; note: string }>("/api/mock-scheduler-tasks"),
        get<{ accounts: Any[] }>("/api/accounts"),
      ])
        .then(([up, local, acc]) => {
          setItems(up.items);
          setMock(local.items);
          setNote(`${up.note} ${local.note}`);
          setAccounts(acc.accounts);
          setDraft((v) =>
            v.account_uid
              ? v
              : acc.accounts[0]
                ? { ...v, account_uid: String(acc.accounts[0].uid) }
                : v,
          );
          setEdits(
            Object.fromEntries(local.items.map((t) => [String(t.id), t])),
          );
        })
        .catch((e) => setError(e.message)),
    [],
  );
  useEffect(load, [load]);
  const inspect = (x: Any) =>
    void get<{ item: Any }>(
      `/api/scheduler-tasks/${encodeURIComponent(String(x.uid))}/${encodeURIComponent(String(x.id))}`,
    )
      .then((d) => setDetail(d.item))
      .catch((e) => setError(e.message));
  const create = async () => {
    try {
      await post("/api/mock-scheduler-tasks", draft);
      setDraft((v) => ({ ...v, name: "", prompt: "" }));
      load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "创建 Mock 任务失败");
    }
  };
  const update = async (id: string) => {
    try {
      const task = edits[id];
      await request(`/api/mock-scheduler-tasks/${encodeURIComponent(id)}`, {
        method: "PUT",
        body: JSON.stringify({
          account_uid: String(task.account_uid || ""),
          name: String(task.name || ""),
          cron: String(task.cron || ""),
          prompt: String(task.prompt || ""),
          enabled: Boolean(task.enabled),
        }),
      });
      load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "保存 Mock 任务失败");
    }
  };
  const remove = async (id: string) => {
    if (!window.confirm("确认删除这条本地 Mock 任务？")) return;
    try {
      await request(`/api/mock-scheduler-tasks/${encodeURIComponent(id)}`, {
        method: "DELETE",
      });
      load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "删除 Mock 任务失败");
    }
  };
  if (error) return <ErrorView text={error} />;
  if (!items || !mock) return <Loading />;
  return (
    <>
      <Head
        eyebrow="ACCOUNT SCHEDULER · MOCK"
        title="账号云端定时任务"
        text="上游 scheduler 未授权时，使用本地 Mock 骨架演练任务管理；不会向 WorkBuddy 创建或修改任何任务。"
        action={
          <button className="black" onClick={load}>
            刷新任务
          </button>
        }
      />
      <div className="notice neutral">{note}</div>
      <section className="card mock-form">
        <h2>创建 Mock 云端任务</h2>
        <div className="mock-fields">
          <label>
            账号
            <select
              disabled={readOnly}
              value={draft.account_uid}
              onChange={(e) =>
                setDraft((v) => ({ ...v, account_uid: e.target.value }))
              }
            >
              {accounts.map((a) => (
                <option key={String(a.uid)} value={String(a.uid)}>
                  {String(a.nickname || a.uid)}
                </option>
              ))}
            </select>
          </label>
          <label>
            任务名称
            <input
              disabled={readOnly}
              value={draft.name}
              placeholder="例如：每日工作摘要"
              onChange={(e) =>
                setDraft((v) => ({ ...v, name: e.target.value }))
              }
            />
          </label>
          <label>
            调度表达式
            <input
              disabled={readOnly}
              value={draft.cron}
              onChange={(e) =>
                setDraft((v) => ({ ...v, cron: e.target.value }))
              }
            />
          </label>
          <label>
            任务内容
            <textarea
              disabled={readOnly}
              value={draft.prompt}
              placeholder="Mock 内容，仅本地保存"
              onChange={(e) =>
                setDraft((v) => ({ ...v, prompt: e.target.value }))
              }
            />
          </label>
          <label className="check">
            <input
              disabled={readOnly}
              type="checkbox"
              checked={draft.enabled}
              onChange={(e) =>
                setDraft((v) => ({ ...v, enabled: e.target.checked }))
              }
            />{" "}
            启用
          </label>
        </div>
        <button
          className="black"
          disabled={readOnly || !draft.account_uid || !draft.name.trim()}
          onClick={() => void create()}
        >
          创建 Mock 任务
        </button>
      </section>
      <section className="card table">
        <h2>本地 Mock 任务</h2>
        <p className="muted">
          仅用于先完成页面与交互骨架，数据保存在控制台自己的 state.json。
        </p>
        <table>
          <thead>
            <tr>
              <th>账号</th>
              <th>名称 / 调度</th>
              <th>内容</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {mock.map((t) => {
              const e = edits[String(t.id)] || t;
              return (
                <tr key={String(t.id)}>
                  <td>{String(t.nickname || t.account_uid)}</td>
                  <td>
                    <input
                      disabled={readOnly}
                      value={String(e.name || "")}
                      onChange={(x) =>
                        setEdits((v) => ({
                          ...v,
                          [String(t.id)]: { ...e, name: x.target.value },
                        }))
                      }
                    />
                    <input
                      disabled={readOnly}
                      value={String(e.cron || "")}
                      onChange={(x) =>
                        setEdits((v) => ({
                          ...v,
                          [String(t.id)]: { ...e, cron: x.target.value },
                        }))
                      }
                    />
                  </td>
                  <td>
                    <textarea
                      disabled={readOnly}
                      value={String(e.prompt || "")}
                      onChange={(x) =>
                        setEdits((v) => ({
                          ...v,
                          [String(t.id)]: { ...e, prompt: x.target.value },
                        }))
                      }
                    />
                  </td>
                  <td>
                    <label className="check">
                      <input
                        disabled={readOnly}
                        type="checkbox"
                        checked={Boolean(e.enabled)}
                        onChange={(x) =>
                          setEdits((v) => ({
                            ...v,
                            [String(t.id)]: { ...e, enabled: x.target.checked },
                          }))
                        }
                      />{" "}
                      {e.enabled ? "启用" : "停用"}
                    </label>
                  </td>
                  <td>
                    <button
                      disabled={readOnly}
                      onClick={() => void update(String(t.id))}
                    >
                      保存
                    </button>
                    <button
                      disabled={readOnly}
                      onClick={() => void remove(String(t.id))}
                    >
                      删除
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        {mock.length === 0 && (
          <Empty text="尚无 Mock 云端任务。请先创建一条本地骨架任务。" />
        )}
      </section>
      <section className="card table">
        <h2>上游真实任务（只读）</h2>
        <table>
          <thead>
            <tr>
              <th>账号</th>
              <th>任务</th>
              <th>状态</th>
              <th>更新时间</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {items.map((x, i) => (
              <tr key={`${String(x.uid)}-${i}`}>
                <td>{String(x.nickname || x.uid)}</td>
                <td>{String(x.name || x.id || "—")}</td>
                <td>
                  {x.error ? (
                    <Badge ok={false}>{String(x.error)}</Badge>
                  ) : (
                    <Badge ok={String(x.status).toLowerCase() !== "failed"}>
                      {String(x.status || "未知")}
                    </Badge>
                  )}
                </td>
                <td>{String(x.updated || "—")}</td>
                <td>
                  {Boolean(x.id) && (
                    <button onClick={() => inspect(x)}>查看详情</button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {items.length === 0 && (
          <Empty text="上游暂未返回任何账号云端定时任务。" />
        )}
      </section>
      {detail && (
        <div className="drawer">
          <button className="close" onClick={() => setDetail(null)}>
            ×
          </button>
          <h2>账号云端任务详情</h2>
          <pre>{JSON.stringify(detail, null, 2)}</pre>
        </div>
      )}
    </>
  );
}

function Activities({
  onNotice,
  readOnly,
}: {
  onNotice: (s: string) => void;
  readOnly: boolean;
}) {
  const [items, setItems] = useState<Any[] | null>(null);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState("");
  const [runState, setRunState] = useState<ActivityRun | null>(null);
  const [showLogs, setShowLogs] = useState(false);
  const [detailUID, setDetailUID] = useState("");
  const load = useCallback(() => {
    setError("");
    void Promise.all([
      get<{ items: Any[] }>("/api/activities"),
      get<{ items: ActivityRun[] }>("/api/activity-runs"),
    ])
      .then(([activities, runs]) => {
        setItems(activities.items);
        setRunState((current) => current || runs.items[0] || null);
      })
      .catch((e) => setError(e.message));
  }, []);
  useEffect(load, [load]);
  useEffect(() => {
    if (!runState || runState.status !== "running") return;
    const timer = window.setInterval(() => {
      void get<{ run: ActivityRun }>(
        `/api/activity-runs/${encodeURIComponent(runState.id)}`,
      )
        .then(({ run }) => {
          setRunState(run);
          if (run.status !== "running") {
            window.clearInterval(timer);
            setBusy("");
            onNotice(run.summary);
            load();
          }
        })
        .catch((e) => {
          window.clearInterval(timer);
          setBusy("");
          onNotice(e instanceof Error ? e.message : "执行状态读取失败");
        });
    }, 1000);
    return () => window.clearInterval(timer);
  }, [runState?.id, runState?.status, load, onNotice]);
  const run = async (uid: string, code = "") => {
    const label = code ? "这项活动" : "该账号全部可自动完成的活动";
    if (!window.confirm(`确认完成${label}并领取已满足的奖励？`)) return;
    const key = `${uid}:${code || "all"}`;
    setBusy(key);
    try {
      const r = await post<{ run: ActivityRun }>(
        `/api/accounts/${encodeURIComponent(uid)}/activities/complete`,
        { task_code: code },
      );
      setRunState(r.run);
      setShowLogs(false);
    } catch (e) {
      onNotice(e instanceof Error ? e.message : "活动任务执行失败");
      setBusy("");
    }
  };
  if (error) return <ErrorView text={error} />;
  if (!items) return <Loading />;
  const accountSummaries = [...new Set(items.map((item) => String(item.uid)))].map(
    (uid) => {
      const tasks = items.filter((item) => String(item.uid) === uid);
      return {
        uid,
        nickname: String(tasks[0]?.nickname || uid),
        tasks,
        total: tasks.filter((item) => !item.error).length,
        completed: tasks.filter(
          (item) =>
            String(item.action) === "done" ||
            ["completed", "claimed"].includes(String(item.status).toLowerCase()),
        ).length,
        automatic: tasks.filter((item) => String(item.action) === "complete")
          .length,
        reward: tasks.reduce(
          (sum, item) => sum + (Number(item.reward) || 0),
          0,
        ),
        error: String(tasks.find((item) => item.error)?.error || ""),
      };
    },
  );
  const detailItems = detailUID
    ? items.filter((item) => String(item.uid) === detailUID)
    : [];
  const detailName = accountSummaries.find((item) => item.uid === detailUID)?.nickname;
  return (
    <>
      <Head
        eyebrow="ACTIVITIES"
        title="活动管理"
        text="读取成长任务，可完成社区已验证的任务并在回读确认后领取奖励。"
        action={
          <button className="black" onClick={load}>
            立即探测
          </button>
        }
      />
      <section className="card activity-batch">
        <div className="section-heading">
          <div>
            <h2>账号活动概览</h2>
            <p className="sub">按账号汇总成长任务和全部奖励；进入明细可单独完成某项任务。</p>
          </div>
          {runState && (
            <button className="secondary-button" onClick={() => setShowLogs((value) => !value)}>
              {showLogs ? "收起执行明细" : "查看执行明细"}
            </button>
          )}
        </div>
        <div className="account-activity-list">
          {accountSummaries.map((account) => (
            <article className="account-activity-row" key={account.uid}>
              <div className="account-identity">
                <b>{account.nickname}</b>
                {account.error ? <Badge ok={false}>{account.error}</Badge> : <span>{account.completed}/{account.total} 项已完成</span>}
              </div>
              <div className="account-metric"><span>活动任务</span><b>{account.total}</b></div>
              <div className="account-metric"><span>全部奖励</span><b>{fmt(account.reward)} 积分</b></div>
              <div className="account-row-actions">
                <button className="secondary-button" onClick={() => setDetailUID((current) => current === account.uid ? "" : account.uid)}>
                  {detailUID === account.uid ? "收起任务明细" : "查看任务明细"}
                </button>
                <button className="primary-button" disabled={readOnly || Boolean(busy) || account.automatic === 0} onClick={() => void run(account.uid)}>
                  {busy === `${account.uid}:all` ? `执行中 ${runState?.completed || 0}/${runState?.total || account.automatic}` : `批量完成任务`}
                </button>
              </div>
            </article>
          ))}
        </div>
        {runState && (
          <div className={`activity-progress ${runState.status}`}>
            <div className="progress-copy">
              <b>{runState.status === "running" ? `正在处理 ${runState.nickname}` : runState.status === "completed" ? "本次执行已完成" : "本次执行未完整结束"}</b>
              <span>{runState.completed} / {runState.total} 项 · 已耗时 {formatSeconds(runState.duration_seconds)}</span>
            </div>
            <progress max={Math.max(runState.total, 1)} value={runState.completed} />
            <div className="progress-stats">
              <span>成功 {runState.succeeded}</span>
              <span>已领取 {runState.already}</span>
              <span>跳过 {runState.skipped}</span>
              <span>失败 {runState.failed}</span>
              <span>积分 +{runState.credit}</span>
              <span>能量 +{runState.energy}</span>
            </div>
            <p>{runState.summary}</p>
          </div>
        )}
        {showLogs && runState && (
          <div className="activity-inline-logs">
            <div className="inline-log-heading">
              <div>
                <h3>活动执行明细</h3>
                <p>{runState.nickname} · {runState.completed}/{runState.total} 项 · {formatSeconds(runState.duration_seconds)}</p>
              </div>
              <button className="secondary-button" onClick={() => setShowLogs(false)}>收起</button>
            </div>
            <ol>
              {runState.logs.map((line, index) => <li key={`${index}-${line}`}>{line}</li>)}
            </ol>
            {runState.logs.length === 0 && <Empty text="执行刚开始，暂时还没有明细。" />}
          </div>
        )}
      </section>
      {detailUID && (
        <section className="card activity-detail-section">
          <div className="section-heading">
            <div>
              <h2>{detailName} 的活动任务</h2>
              <p className="sub">共 {detailItems.length} 项，可在这里查看进度或单独执行。</p>
            </div>
            <button className="secondary-button" onClick={() => setDetailUID("")}>关闭明细</button>
          </div>
          <div className="table activity-detail-table"><table>
          <thead>
            <tr>
              <th>账号</th>
              <th>活动任务</th>
              <th>状态</th>
              <th>奖励</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {detailItems.map((x, i) => (
              <tr key={`${String(x.uid)}-${String(x.code)}-${i}`}>
                <td>{String(x.nickname || x.uid)}</td>
                {x.error ? (
                  <td colSpan={4}>
                    <Badge ok={false}>{String(x.error)}</Badge>
                  </td>
                ) : (
                  <>
                    <td>{String(x.name || x.code)}</td>
                    <td>
                      {activityStatus(String(x.status || ""))}
                      {Number(x.target) > 0 && (
                        <small className="sub"> {fmt(x.progress)} / {fmt(x.target)}</small>
                      )}
                    </td>
                    <td>{fmt(x.reward)}</td>
                    <td>
                      {Boolean(x.new) && <Badge ok>新发现</Badge>}
                      {String(x.action) === "complete" && (
                        <button
                          disabled={readOnly || Boolean(busy)}
                          onClick={() => void run(String(x.uid), String(x.code))}
                        >
                          {busy === `${String(x.uid)}:${String(x.code)}` ? "执行中…" : "完成并领奖"}
                        </button>
                      )}
                      {String(x.action) === "done" && <Badge ok>已领取</Badge>}
                      {String(x.action) === "manual" && <span className="sub">需人工完成</span>}
                    </td>
                  </>
                )}
              </tr>
            ))}
          </tbody>
          </table></div>
          {detailItems.length === 0 && <Empty text="当前账号没有可读取的成长任务。" />}
        </section>
      )}
    </>
  );
}

function formatSeconds(value: number) {
  const seconds = Math.max(0, Math.floor(value || 0));
  if (seconds < 60) return `${seconds} 秒`;
  return `${Math.floor(seconds / 60)} 分 ${seconds % 60} 秒`;
}

function activityStatus(value: string) {
  const labels: Record<string, string> = {
    not_accepted: "未接受",
    accepted: "已接受",
    in_progress: "进行中",
    completed: "已完成，待领取",
    claimed: "已领取",
    failed: "执行失败",
  };
  return labels[value.toLowerCase()] || value || "未知";
}

function OAuth({
  onNotice,
  onDone,
  readOnly,
}: {
  onNotice: (s: string) => void;
  onDone: () => void;
  readOnly: boolean;
}) {
  const [region, setRegion] = useState("cn");
  const [flow, setFlow] = useState<{ id: string; url: string } | null>(null);
  const [starting, setStarting] = useState(false);
  const [status, setStatus] = useState("");
  const start = async () => {
    setStarting(true);
    setStatus("");
    try {
      setFlow(
        (await post("/api/oauth/start", { region })) as {
          id: string;
          url: string;
        },
      );
      setStatus("等待你在 WorkBuddy 页面完成授权…");
    } catch (e) {
      onNotice(e instanceof Error ? e.message : "发起失败");
    } finally {
      setStarting(false);
    }
  };
  useEffect(() => {
    if (!flow) return;
    let failures = 0;
    const timer = window.setInterval(
      () =>
        void post<{ status: string; nickname?: string }>(
          `/api/oauth/${encodeURIComponent(flow.id)}/poll`,
        )
          .then((r) => {
            failures = 0;
            if (r.status === "success") {
              window.clearInterval(timer);
              onNotice(`账号 ${r.nickname || ""} 已保存`);
              onDone();
            }
          })
          .catch((e) => {
            failures++;
            if (failures >= 3) {
              window.clearInterval(timer);
              setStatus(
                e instanceof Error
                  ? `授权轮询已停止：${e.message}`
                  : "授权轮询已停止，请重新发起",
              );
            }
          }),
      3000,
    );
    return () => window.clearInterval(timer);
  }, [flow, onDone, onNotice]);
  return (
    <>
      <Head
        eyebrow="OAUTH"
        title="添加账号"
        text="选择账号区域后完成 WorkBuddy 网页授权；凭据只写入共享 auths 目录。"
      />
      <section className="card oauth">
        <label>
          账号区域
          <select
            disabled={readOnly}
            value={region}
            onChange={(e) => setRegion(e.target.value)}
          >
            <option value="cn">中国大陆</option>
            <option value="global">国际版</option>
          </select>
        </label>
        <button
          className="black"
          disabled={readOnly || starting}
          onClick={() => void start()}
        >
          {starting ? "正在获取授权链接…" : "发起 OAuth 授权"}
        </button>
        {status && (
          <p className="muted" role="status">
            {status}
          </p>
        )}
        {flow && (
          <div className="oauth-link">
            <p>请在新窗口完成授权，控制台会自动轮询结果。</p>
            <a
              className="black"
              href={flow.url}
              target="_blank"
              rel="noreferrer"
            >
              打开 WorkBuddy 授权页
            </a>
          </div>
        )}
      </section>
    </>
  );
}

function Badge({ ok, children }: { ok?: boolean; children: React.ReactNode }) {
  return (
    <span className={ok === false ? "badge bad" : "badge"}>{children}</span>
  );
}
function Empty({ text }: { text: string }) {
  return <div className="empty">{text}</div>;
}
