const categories = [
  "Tất cả",
  "Sân 5",
  "Sân 7",
  "Sân 11",
  "Ngoài trời",
  "Có mái che",
  "Đá ban đêm",
  "Trẻ em",
];

const featuredFields = [
  {
    name: "Sân bóng mini Hưng Thịnh",
    type: "Sân 5",
    location: "Quận 7, TP.HCM",
    price: "250k / giờ",
    image:
      "https://images.unsplash.com/photo-1574629810360-7efbbe195018?auto=format&fit=crop&w=900&q=80",
    badge: "Mới",
  },
  {
    name: "Sân bóng Tân Phú",
    type: "Sân 7",
    location: "Quận Tân Phú",
    price: "350k / giờ",
    image:
      "https://images.unsplash.com/photo-1547347298-4074fc3086f0?auto=format&fit=crop&w=900&q=80",
    badge: "Hot",
  },
  {
    name: "Sân bóng Thủ Đức",
    type: "Sân 11",
    location: "Thủ Đức",
    price: "550k / giờ",
    image:
      "https://images.unsplash.com/photo-1517649763962-0c623066013b?auto=format&fit=crop&w=900&q=80",
    badge: "Phổ biến",
  },
  {
    name: "Sân bóng Cầu Giấy",
    type: "Sân 5",
    location: "Cầu Giấy",
    price: "280k / giờ",
    image:
      "https://images.unsplash.com/photo-1579952363873-27d3bfad9c0d?auto=format&fit=crop&w=900&q=80",
    badge: "Mới",
  },
];

const steps = [
  {
    number: "1",
    title: "Tìm sân",
    text: "Lọc theo khu vực, loại sân, thời gian và giá phù hợp.",
  },
  {
    number: "2",
    title: "Chọn giờ",
    text: "Chọn khung giờ rảnh và xem thông tin sân trong từng khung giờ.",
  },
  {
    number: "3",
    title: "Thanh toán",
    text: "Xác nhận đặt sân nhanh, tiện dụng và an toàn với nhiều phương thức thanh toán.",
  },
];

const benefits = [
  "Đặt sân trong 30 giây",
  "Nhiều sân chất lượng gần bạn",
  "Giá minh bạch, không phí ẩn",
  "Hỗ trợ 24/7",
  "Thanh toán linh hoạt",
  "Đảm bảo vệ sinh, ánh sáng",
];

const testimonials = [
  {
    name: "Anh Quốc",
    role: "Đội bóng mini",
    rating: 5,
    text: "Ứng dụng dễ dùng, đặt sân nhanh, sân sạch và thiết bị đầy đủ. Chúng tôi rất hài lòng.",
  },
  {
    name: "Chị Lan",
    role: "Nhóm bạn thân",
    rating: 5,
    text: "Mình đặt sân qua app ở khu vực gần nhà. Giá rõ ràng, lịch đặt trực quan, rất tiện.",
  },
  {
    name: "Anh Dũng",
    role: "Huấn luyện viên",
    rating: 5,
    text: "Sân phong phú, không gian rộng, có cả sân đêm. Hỗ trợ tốt và thanh toán rất dễ.",
  },
];

const footerLinks = {
  "Sản phẩm": ["Đặt sân", "Quản lý sân", "Tích hợp API"],
  "Hỗ trợ": ["Điều khoản", "Chính sách", "Liên hệ"],
  "Khu vực": ["TP.HCM", "Hà Nội", "Đà Nẵng"],
};

export default function Home() {
  return (
    <main className="min-h-screen bg-[#f3f3f3] text-slate-900">
      <header className="bg-[radial-gradient(circle_at_top_left,_rgba(52,211,153,0.28),_transparent_30%),linear-gradient(135deg,#0f4f43_0%,#0a3c31_35%,#0f5b4d_100%)] text-white">
        <div className="mx-auto max-w-6xl px-4 pb-12 pt-5 sm:px-6 lg:px-8">
          <nav className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-2xl border border-white/20 bg-white/10 text-lg font-bold backdrop-blur-sm">
                SB
              </div>
              <div>
                <div className="text-xl font-black tracking-tight">SânBóng.vn</div>
              </div>
            </div>

            <div className="hidden items-center gap-8 text-sm text-emerald-50/90 md:flex">
              <a href="#">Trang chủ</a>
              <a href="#">Sân bóng</a>
              <a href="#">Dịch vụ</a>
              <a href="#">Tin tức</a>
              <a href="#">Liên hệ</a>
            </div>

            <button className="rounded-full bg-[#69d99e] px-5 py-2.5 text-sm font-semibold text-[#083d31] shadow-lg shadow-emerald-900/20 transition hover:bg-[#84e7b0]">
              Đặt sân ngay
            </button>
          </nav>

          <div className="mt-10 grid items-center gap-8 lg:grid-cols-[1.2fr_0.8fr]">
            <div>
              <div className="mb-5 inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/5 px-3 py-1.5 text-xs font-medium text-emerald-100 backdrop-blur-sm">
                <span className="inline-block h-2 w-2 rounded-full bg-[#8ef0b0]" />
                Đặt sân nhanh chóng, tiện lợi
              </div>

              <h1 className="max-w-xl text-4xl font-black leading-tight sm:text-5xl">
                Đặt sân bóng địa phương dễ dàng hơn bao giờ hết
              </h1>

              <p className="mt-4 max-w-xl text-base text-emerald-50/85 sm:text-lg">
                Hàng trăm sân bóng chất lượng, giá tốt, dễ tìm theo khu vực và thời gian phù hợp với bạn.
              </p>

              <div className="mt-7 flex flex-wrap items-center gap-5">
                <div className="flex items-center gap-2 rounded-full bg-white/8 px-3 py-2 text-sm text-emerald-50">
                  <span className="text-lg">⭐</span>
                  4.9/5 đánh giá từ khách hàng
                </div>
                <div className="text-sm text-emerald-100">3250+ lượt đặt trong tháng</div>
              </div>

              <div className="mt-8 flex flex-wrap gap-4 text-sm text-emerald-50/80">
                <div className="rounded-full border border-white/15 bg-white/5 px-4 py-2">623 sân bóng</div>
                <div className="rounded-full border border-white/15 bg-white/5 px-4 py-2">120 khu vực</div>
                <div className="rounded-full border border-white/15 bg-white/5 px-4 py-2">24/7 hỗ trợ</div>
              </div>
            </div>

            <div className="rounded-[30px] border border-white/10 bg-white/8 p-4 shadow-2xl shadow-emerald-950/30 backdrop-blur-sm">
              <div className="rounded-[24px] bg-[#e8f7ee] p-5 text-slate-800">
                <div className="mb-4 flex items-center justify-between">
                  <div>
                    <p className="text-xs font-medium uppercase tracking-[0.2em] text-emerald-700">Đặt sân</p>
                    <h2 className="mt-1 text-2xl font-black text-slate-900">Tìm sân ngay</h2>
                  </div>
                  <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#d4f7de] text-lg">⚽</div>
                </div>

                <div className="space-y-3 text-sm">
                  <div className="rounded-2xl bg-white px-4 py-3 shadow-sm">
                    <div className="text-[11px] font-bold uppercase tracking-[0.2em] text-slate-500">Khu vực</div>
                    <div className="mt-1 font-semibold">TP.HCM</div>
                  </div>
                  <div className="rounded-2xl bg-white px-4 py-3 shadow-sm">
                    <div className="text-[11px] font-bold uppercase tracking-[0.2em] text-slate-500">Loại sân</div>
                    <div className="mt-1 font-semibold">Sân 5, 7, 11</div>
                  </div>
                  <div className="rounded-2xl bg-white px-4 py-3 shadow-sm">
                    <div className="text-[11px] font-bold uppercase tracking-[0.2em] text-slate-500">Thời gian</div>
                    <div className="mt-1 font-semibold">Hôm nay, 19:00 - 21:00</div>
                  </div>
                </div>

                <button className="mt-5 w-full rounded-2xl bg-[#0e5b4a] px-4 py-3 text-base font-bold text-white shadow-lg shadow-emerald-900/20 transition hover:bg-[#0c4d3d]">
                  Tìm sân ngay
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="mx-auto -mt-8 max-w-6xl px-4 sm:px-6 lg:px-8">
        <section className="rounded-[28px] bg-white p-6 shadow-[0_18px_50px_rgba(15,36,29,0.08)] ring-1 ring-slate-200/70">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700">Dịch vụ</p>
              <h3 className="mt-1 text-2xl font-black text-slate-900">Chọn sân phù hợp</h3>
            </div>

            <div className="flex flex-wrap gap-2">
              {categories.map((category, index) => (
                <button
                  key={category}
                  className={`rounded-full px-4 py-2 text-sm font-medium transition ${
                    index === 0
                      ? "bg-[#0f5b4d] text-white shadow-md shadow-emerald-900/10"
                      : "bg-[#eef7f1] text-slate-700 hover:bg-[#dff3e9]"
                  }`}
                >
                  {category}
                </button>
              ))}
            </div>
          </div>
        </section>
      </div>

      <main className="mx-auto max-w-6xl space-y-12 px-4 pb-16 pt-10 sm:px-6 lg:px-8">
        <section>
          <div className="mb-6 flex items-center justify-between gap-3">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700">Nổi bật</p>
              <h3 className="mt-2 text-3xl font-black text-slate-900">Các sân đang được yêu thích</h3>
            </div>
            <button className="hidden rounded-full border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 transition hover:border-emerald-200 hover:text-emerald-700 sm:inline-flex">
              Xem tất cả
            </button>
          </div>

          <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
            {featuredFields.map((field) => (
              <article key={field.name} className="overflow-hidden rounded-[28px] bg-white shadow-[0_16px_30px_rgba(15,23,42,0.08)] ring-1 ring-slate-100">
                <div className="relative">
                  <img src={field.image} alt={field.name} className="h-52 w-full object-cover" />
                  <span className="absolute left-4 top-4 rounded-full bg-[#0f5b4d] px-3 py-1 text-xs font-bold text-white shadow-md">
                    {field.badge}
                  </span>
                </div>
                <div className="p-4">
                  <div className="mb-2 flex items-center justify-between gap-2">
                    <span className="rounded-full bg-[#e3f6ea] px-2.5 py-1 text-[11px] font-bold uppercase tracking-[0.15em] text-emerald-700">
                      {field.type}
                    </span>
                    <span className="text-lg font-black text-[#0d5c4c]">{field.price}</span>
                  </div>
                  <h4 className="text-lg font-bold text-slate-900">{field.name}</h4>
                  <p className="mt-2 text-sm text-slate-500">📍 {field.location}</p>
                  <div className="mt-4 flex items-center justify-between">
                    <div className="flex items-center gap-1 text-amber-400">★★★★★</div>
                    <button className="rounded-full bg-[#0f5b4d] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#0d4d40]">
                      Đặt ngay
                    </button>
                  </div>
                </div>
              </article>
            ))}
          </div>
        </section>

        <section className="rounded-[30px] bg-[#eef7f1] p-6 sm:p-8">
          <div className="mb-8 text-center">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700">Quy trình</p>
            <h3 className="mt-2 text-3xl font-black text-slate-900">Đặt sân trong 3 bước</h3>
          </div>

          <div className="grid gap-6 md:grid-cols-3">
            {steps.map((step, index) => (
              <div key={step.number} className="rounded-[28px] bg-white p-5 shadow-sm ring-1 ring-slate-100">
                <div className="mb-5 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#d9f5e7] text-lg font-black text-[#0f5b4d]">
                  {step.number}
                </div>
                <div className="mb-2 inline-flex rounded-full bg-[#eefbf3] px-3 py-1 text-xs font-bold uppercase tracking-[0.12em] text-emerald-700">
                  Bước {index + 1}
                </div>
                <h4 className="text-xl font-bold text-slate-900">{step.title}</h4>
                <p className="mt-3 text-sm leading-6 text-slate-600">{step.text}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="grid gap-8 lg:grid-cols-[1.1fr_0.9fr] lg:items-center">
          <div className="rounded-[32px] bg-white p-5 shadow-[0_18px_50px_rgba(15,36,29,0.06)] ring-1 ring-slate-100 sm:p-7">
            <div className="mb-6 flex items-center justify-between gap-3">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700">Ưu điểm</p>
                <h3 className="mt-2 text-3xl font-black text-slate-900">Đặt sân nhanh, an toàn</h3>
              </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              {benefits.map((item) => (
                <div key={item} className="flex items-center gap-3 rounded-2xl bg-[#f7faf8] p-4 text-sm font-medium text-slate-700">
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-[#cfeee0] text-base text-[#0d5c4c]">✓</div>
                  {item}
                </div>
              ))}
            </div>
          </div>

          <div className="overflow-hidden rounded-[32px] bg-[linear-gradient(180deg,#0b3d31_0%,#134b3f_100%)] p-6 text-white shadow-[0_20px_50px_rgba(12,57,50,0.35)]">
            <div className="mb-5 flex items-center justify-between">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-200">Sân gần bạn</p>
                <h3 className="mt-2 text-3xl font-black">Sân bóng giao lưu</h3>
              </div>
              <button className="rounded-full bg-white/10 px-3 py-2 text-xs font-semibold text-emerald-50 transition hover:bg-white/15">
                Xem ngay
              </button>
            </div>

            <div className="rounded-[24px] bg-white/6 p-3 backdrop-blur-sm">
              <img
                src="https://images.unsplash.com/photo-1521412644187-c49fa049e84d?auto=format&fit=crop&w=1200&q=80"
                alt="Sân bóng"
                className="h-72 w-full rounded-[20px] object-cover"
              />
            </div>

            <div className="mt-5 grid gap-3 sm:grid-cols-3">
              <div className="rounded-2xl bg-white/7 p-3">
                <div className="text-2xl font-black">120+</div>
                <div className="mt-1 text-xs uppercase tracking-[0.2em] text-emerald-100">Sân</div>
              </div>
              <div className="rounded-2xl bg-white/7 p-3">
                <div className="text-2xl font-black">3.2k</div>
                <div className="mt-1 text-xs uppercase tracking-[0.2em] text-emerald-100">Lượt chơi</div>
              </div>
              <div className="rounded-2xl bg-white/7 p-3">
                <div className="text-2xl font-black">24/7</div>
                <div className="mt-1 text-xs uppercase tracking-[0.2em] text-emerald-100">Hỗ trợ</div>
              </div>
            </div>
          </div>
        </section>

        <section className="space-y-8">
          <div className="text-center">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700">Đánh giá</p>
            <h3 className="mt-2 text-3xl font-black text-slate-900">Khách hàng nói gì về chúng tôi</h3>
          </div>

          <div className="grid gap-6 lg:grid-cols-3">
            {testimonials.map((item) => (
              <div key={item.name} className="rounded-[28px] bg-white p-6 shadow-[0_16px_30px_rgba(15,23,42,0.06)] ring-1 ring-slate-100">
                <div className="mb-4 flex items-center gap-1 text-lg text-amber-400">
                  {Array.from({ length: item.rating }).map((_, idx) => (
                    <span key={`${item.name}-${idx}`}>★</span>
                  ))}
                </div>
                <p className="text-sm leading-7 text-slate-600">“{item.text}”</p>
                <div className="mt-6 flex items-center gap-3">
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-[#d8f4e5] text-sm font-black text-[#0f5b4d]">
                    {item.name.slice(0, 1)}
                  </div>
                  <div>
                    <div className="font-bold text-slate-900">{item.name}</div>
                    <div className="text-xs uppercase tracking-[0.12em] text-slate-500">{item.role}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>

        <section className="rounded-[32px] bg-[#0f2f29] px-5 py-14 text-white shadow-[0_25px_60px_rgba(11,61,49,0.25)] sm:px-8 lg:px-12">
          <div className="grid gap-8 lg:grid-cols-[1.1fr_0.9fr] lg:items-center">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.25em] text-emerald-200">Phần mềm</p>
              <h3 className="mt-3 text-4xl font-black leading-tight">Phần mềm quản lý sân bóng thông minh</h3>
              <ul className="mt-6 space-y-3 text-sm text-emerald-50/80">
                <li>• Theo dõi lịch đặt, giờ hoạt động và doanh thu rõ ràng</li>
                <li>• Cập nhật sân, khuyến mãi và thanh toán online dễ dàng</li>
                <li>• Quản lý nhân sự, khách hàng và báo cáo nhanh chóng</li>
              </ul>
              <button className="mt-7 rounded-full bg-[#69d99e] px-5 py-3 text-base font-bold text-[#083d31] transition hover:bg-[#84e7b0]">
                Tải app ngay
              </button>
            </div>

            <div className="rounded-[28px] bg-white/7 p-5 backdrop-blur-sm">
              <div className="rounded-[24px] bg-[#eafaf1] p-5 text-slate-900">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[#0f5b4d] text-sm font-black text-white">SB</div>
                    <div>
                      <div className="text-sm font-bold text-slate-500">SânBóng</div>
                      <div className="text-lg font-black">Quản lý sân</div>
                    </div>
                  </div>
                  <span className="rounded-full bg-[#d7f7e6] px-2.5 py-1 text-[10px] font-bold uppercase tracking-[0.15em] text-emerald-700">
                    Online
                  </span>
                </div>

                <div className="mt-5 grid gap-3 sm:grid-cols-2">
                  <div className="rounded-2xl bg-white p-4 shadow-sm">
                    <div className="text-[11px] font-bold uppercase tracking-[0.22em] text-slate-500">Doanh thu</div>
                    <div className="mt-2 text-2xl font-black text-[#0b3d31]">34.8M</div>
                  </div>
                  <div className="rounded-2xl bg-white p-4 shadow-sm">
                    <div className="text-[11px] font-bold uppercase tracking-[0.22em] text-slate-500">Đặt sân</div>
                    <div className="mt-2 text-2xl font-black text-[#0b3d31]">1,260</div>
                  </div>
                </div>

                <div className="mt-4 rounded-2xl bg-[#f2fbf5] p-4">
                  <div className="mb-3 flex items-center justify-between text-sm font-semibold text-slate-600">
                    <span>Khung giờ nổi bật</span>
                    <span>7:00 - 22:00</span>
                  </div>
                  <div className="space-y-2">
                    {[60, 80, 92, 68].map((value, idx) => (
                      <div key={idx} className="flex items-center gap-2">
                        <span className="w-9 text-xs text-slate-500">{idx + 1}</span>
                        <div className="h-2.5 flex-1 overflow-hidden rounded-full bg-[#e3f1ea]">
                          <div className="h-full rounded-full bg-[#0f5b4d]" style={{ width: `${value}%` }} />
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>

      <footer className="bg-[#081b18] text-slate-200">
        <div className="mx-auto max-w-6xl px-4 py-14 sm:px-6 lg:px-8">
          <div className="grid gap-9 md:grid-cols-[1.3fr_0.8fr_0.8fr_0.8fr]">
            <div>
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#0f5b4d] font-black text-white">SB</div>
                <div>
                  <div className="text-xl font-black text-white">SânBóng.vn</div>
                </div>
              </div>
              <p className="mt-4 max-w-sm text-sm leading-7 text-slate-300">
                Nền tảng đặt sân bóng địa phương, giúp bạn tìm sân nhanh, tiết kiệm thời gian và dễ quản lý hoạt động thể thao.
              </p>
            </div>

            {Object.entries(footerLinks).map(([title, items]) => (
              <div key={title}>
                <h4 className="text-lg font-bold text-white">{title}</h4>
                <ul className="mt-4 space-y-3 text-sm text-slate-300">
                  {items.map((item) => (
                    <li key={item}>{item}</li>
                  ))}
                </ul>
              </div>
            ))}
          </div>

          <div className="mt-10 flex flex-col gap-3 border-t border-white/10 pt-6 text-sm text-slate-400 sm:flex-row sm:items-center sm:justify-between">
            <span>© 2026 SânBóng.vn. All rights reserved.</span>
            <div className="flex gap-4">
              <span>Facebook</span>
              <span>Instagram</span>
              <span>Zalo</span>
            </div>
          </div>
        </div>
      </footer>
    </main>
  );
}
