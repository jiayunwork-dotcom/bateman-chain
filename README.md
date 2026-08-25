# bateman-chain

bateman-chain 是放射性线性衰变链求解器：给定各核素衰变常数 λ（1/s）与初始核数，命令行沿时间输出每一核的核数 N 与活度 A（A = λ·N）。链方程 dNi/dt = λ_{i-1}·N_{i-1} − λ_i·Ni，末核可稳定（λk = 0）或继续衰变；闭式 Bateman 求和与「原子不灭、只在链上迁移」一致：稳定终点的核数 + 仍在链上的核数 = 初始总原子。

- 输入：JSON 规格文件（`lambda` 衰变常数数组、`initial` 初始核数数组；可选 `names` 核素标签），示例见 `example/u238-short.json`
- 输出：`eval` 按 `--t`（单点）、`--times`（时刻列表）或 `--t0/--t1/--steps`（均匀网格）打印各核 N 与 A、链上总数与已衰变原子数；`--format csv` 输出 CSV
- 边界：链长 2～12；所有 λ 非负且有限；初始核数非负且有限；时间非负；输出点数与数值积分步数均有上限，超限返回 error
- 非法输入（λ<0、初始核数为负、链长 <2 或 >12、t<0、网格点数超限、积分步数预算超限）一律返回 error，CLI 打印 stderr 并非零退出

## 钉死的约定

- **方程**：线性链 N1→N2→…→Nk，`dNi/dt = λ_{i-1}·N_{i-1} − λ_i·Ni`；末核 λk = 0 表示稳定（长期核数收敛到初始总原子），λk > 0 表示继续衰变（长期全部归零）。
- **闭式解**：优先 Bateman 求和。对每个核 i，`Ni(t) = Σ_{j≤i} Nj(0)·(Π_{l=j}^{i-1} λ_l)·Σ_{p=j}^{i} e^{−λ_p·t}/Π_{q≠p}(λ_q−λ_p)`；活度与闭式系数共用同一 λ 向量，`A_i = λ_i·N_i`。
- **避免除零**：Bateman 分母是 λ_q − λ_p 的乘积。当任一对 λ 的相对差 ≤ `1e-8`（`internal/bateman` 的 `NearTol`）时，闭式分母会逼近零、除法放大误差；此时求解自动改走 `internal/integrate` 的 RK4 数值积分，全程不对任何 `λ_q − λ_p` 做除法，因此不会除零。t=0 恒直接返回初值，也不涉及除法。
- **步数上限**：`--steps` 网格点数受 `--steps-limit`（默认 100000）约束；RK4 路径还受积分步数预算（默认 1000000，按最快衰变率自适应取步长）约束；超限返回 error。
- **守恒与长期极限**：全部 λ>0 时链上总核数严格下降，母体开端下降率 = λ1·N1(0)；末核稳定时 N_end → Σ Ni(0)，全不稳定时全部 → 0；λ_daughter ≫ λ_parent 时短寿命子体长期逼近 λ_parent·N_parent/λ_daughter。
- **t=0**：两种方法都原样返回初始核数，活度 = λ∘N0。

## example/u238-short.json 期望

三核链 U-238 → Th-234 → Pa-234m，λ = [4.916e-18, 2.82e-6, 2.12e-3]，初始仅 U-238 有 1e6 个原子。时间短到母体几乎不变：`λ1·t = 1.77e-14`，U-238 仍 ≈ 1e6；子体按 `λ1·N1·t ≈ 1.77e-8` 量级生长（3600 s 时 Th-234 ≈ 1.76e-8），Pa-234m 处于短寿命准平衡 ≈ λ2·N2/λ3 ≈ 2e-11。

## 构建 / 运行 / 测试

```text
go build ./...                     # 编译（纯标准库）
go test ./...                      # 全部测试（chain / bateman / integrate / eval）
go run . eval example/u238-short.json --t 3600
go run . eval example/u238-short.json --t 3600 --format csv
go run . eval example/u238-short.json --times 0,600,1800,3600
go run . eval example/u238-short.json --t0 0 --t1 3600 --steps 5
```

非法输入抽查（均 stderr + 非零退出码）：

```text
go run . eval example/u238-short.json --t -5                 # time must be non-negative
go run . eval /tmp/one-nuclide.json --t 10                   # chain must have at least 2 nuclides
go run . eval /tmp/negative-lambda.json --t 10               # lambda must be non-negative
```
