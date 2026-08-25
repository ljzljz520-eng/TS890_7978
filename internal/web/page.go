package web

const pageHTML = `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>宠物短片色彩调试页</title>
<style>
:root {
  color-scheme: dark;
  --ink: #f5f7f8;
  --muted: #9ca7ad;
  --surface: #171b1d;
  --surface-raised: #202629;
  --surface-soft: #293034;
  --line: #364044;
  --accent: #43d17d;
  --accent-strong: #20b968;
  --warning: #ffbf47;
  --danger: #ff6b6b;
  --focus: #68b5ff;
  --radius: 6px;
  --sidebar: 264px;
  --inspector: 328px;
}
* {
  box-sizing: border-box;
}
html,
body {
  width: 100%;
  min-height: 100%;
  margin: 0;
  background: #0f1213;
  color: var(--ink);
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  font-size: 14px;
  letter-spacing: 0;
}
button,
input,
select {
  font: inherit;
}
button {
  border: 0;
}
button:focus-visible,
input:focus-visible,
select:focus-visible {
  outline: 2px solid var(--focus);
  outline-offset: 2px;
}
.app {
  display: grid;
  grid-template-columns: var(--sidebar) minmax(0, 1fr) var(--inspector);
  grid-template-rows: 56px minmax(0, 1fr) 96px;
  min-height: 100vh;
}
.topbar {
  grid-column: 1 / -1;
  display: flex;
  align-items: center;
  gap: 16px;
  min-width: 0;
  padding: 0 18px;
  border-bottom: 1px solid var(--line);
  background: #141819;
}
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  width: calc(var(--sidebar) - 18px);
  flex: 0 0 auto;
}
.brand-mark {
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  border-radius: 6px;
  background: var(--accent);
  color: #08150d;
  font-weight: 800;
}
.brand-copy {
  min-width: 0;
}
.brand-title {
  margin: 0;
  font-size: 14px;
  line-height: 18px;
  font-weight: 700;
}
.brand-subtitle {
  color: var(--muted);
  font-size: 11px;
  line-height: 14px;
}
.project-name {
  min-width: 0;
  overflow: hidden;
  color: #cbd2d5;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.top-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}
.icon-button,
.command-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 34px;
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: var(--surface-raised);
  color: var(--ink);
  cursor: pointer;
}
.icon-button {
  width: 34px;
  padding: 0;
}
.command-button {
  gap: 8px;
  padding: 0 12px;
  font-weight: 650;
}
.command-button.primary {
  border-color: var(--accent);
  background: var(--accent);
  color: #07130b;
}
.command-button:hover,
.icon-button:hover {
  filter: brightness(1.12);
}
.sidebar {
  grid-column: 1;
  grid-row: 2 / 4;
  min-height: 0;
  border-right: 1px solid var(--line);
  background: var(--surface);
  overflow: auto;
}
.sidebar-section {
  padding: 16px 14px;
  border-bottom: 1px solid var(--line);
}
.section-label {
  margin: 0 0 10px;
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
}
.upload-zone {
  display: grid;
  place-items: center;
  min-height: 120px;
  padding: 16px;
  border: 1px dashed #536065;
  border-radius: var(--radius);
  background: #1b2022;
  color: #c4cccf;
  text-align: center;
  cursor: pointer;
}
.upload-zone:hover {
  border-color: var(--accent);
  background: #1d2621;
}
.upload-icon {
  display: grid;
  place-items: center;
  width: 38px;
  height: 38px;
  margin-bottom: 8px;
  border-radius: 50%;
  background: var(--surface-soft);
  font-size: 21px;
}
.upload-main {
  font-weight: 650;
}
.upload-note {
  margin-top: 4px;
  color: var(--muted);
  font-size: 11px;
  line-height: 16px;
}
.filter-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}
.field-label {
  display: block;
  margin-bottom: 6px;
  color: #bec6c9;
  font-size: 12px;
}
.select,
.search {
  width: 100%;
  height: 34px;
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: #101415;
  color: var(--ink);
  padding: 0 9px;
}
.search {
  margin-bottom: 8px;
}
.clip-list {
  display: grid;
  gap: 6px;
}
.clip-item {
  display: grid;
  grid-template-columns: 46px minmax(0, 1fr);
  gap: 9px;
  padding: 8px;
  border: 1px solid transparent;
  border-radius: var(--radius);
  background: transparent;
  color: inherit;
  text-align: left;
  cursor: pointer;
}
.clip-item:hover {
  background: var(--surface-raised);
}
.clip-item.active {
  border-color: #447d5c;
  background: #1f2d25;
}
.clip-thumb {
  display: grid;
  place-items: center;
  width: 46px;
  height: 34px;
  border-radius: 4px;
  background: #30383c;
  color: #d9e0e2;
  font-size: 17px;
}
.clip-meta {
  min-width: 0;
}
.clip-name {
  overflow: hidden;
  font-size: 12px;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.clip-detail {
  margin-top: 3px;
  color: var(--muted);
  font-size: 11px;
}
.workspace {
  grid-column: 2;
  grid-row: 2;
  display: grid;
  grid-template-rows: 42px minmax(0, 1fr);
  min-width: 0;
  min-height: 0;
  background: #0d1011;
}
.workspace-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  border-bottom: 1px solid var(--line);
  background: #171b1d;
}
.segmented {
  display: inline-flex;
  height: 28px;
  padding: 2px;
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: #101415;
}
.segment {
  min-width: 66px;
  border-radius: 4px;
  background: transparent;
  color: var(--muted);
  cursor: pointer;
}
.segment.active {
  background: var(--surface-soft);
  color: var(--ink);
}
.zoom-copy {
  margin-left: auto;
  color: var(--muted);
  font-size: 12px;
}
.stage {
  display: grid;
  place-items: center;
  min-width: 0;
  min-height: 0;
  padding: 24px;
  overflow: auto;
}
.viewer {
  position: relative;
  width: min(100%, 960px);
  aspect-ratio: 16 / 9;
  border: 1px solid #414b4f;
  background:
    linear-gradient(45deg, #171b1d 25%, transparent 25%),
    linear-gradient(-45deg, #171b1d 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #171b1d 75%),
    linear-gradient(-45deg, transparent 75%, #171b1d 75%),
    #111516;
  background-position: 0 0, 0 8px, 8px -8px, -8px 0;
  background-size: 16px 16px;
  box-shadow: 0 16px 40px rgb(0 0 0 / 35%);
  overflow: hidden;
}
.viewer-empty {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  color: var(--muted);
  text-align: center;
}
.viewer-empty strong {
  display: block;
  margin-bottom: 6px;
  color: #d9dfe1;
  font-size: 16px;
}
.viewer-photo {
  position: absolute;
  inset: 0;
  display: none;
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.viewer-overlay {
  position: absolute;
  right: 12px;
  bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: 4px;
  background: rgb(0 0 0 / 62%);
  color: #e7ebec;
  font-size: 11px;
}
.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--accent);
}
.inspector {
  grid-column: 3;
  grid-row: 2 / 4;
  min-height: 0;
  border-left: 1px solid var(--line);
  background: var(--surface);
  overflow: auto;
}
.inspector-header {
  position: sticky;
  top: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  height: 42px;
  padding: 0 14px;
  border-bottom: 1px solid var(--line);
  background: var(--surface);
  font-size: 13px;
  font-weight: 700;
}
.control-section {
  padding: 16px 14px;
  border-bottom: 1px solid var(--line);
}
.preset-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 7px;
}
.preset-button {
  min-height: 62px;
  padding: 7px 5px;
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: var(--surface-raised);
  color: var(--muted);
  cursor: pointer;
}
.preset-button strong {
  display: block;
  margin-top: 4px;
  color: var(--ink);
  font-size: 11px;
}
.preset-button.active {
  border-color: var(--accent);
  background: #203228;
}
.slider-control {
  display: grid;
  grid-template-columns: 74px minmax(0, 1fr) 42px;
  align-items: center;
  gap: 8px;
  min-height: 38px;
}
.slider-control label {
  color: #cbd2d5;
  font-size: 12px;
}
.slider-control input[type="range"] {
  width: 100%;
  accent-color: var(--accent);
}
.slider-value {
  width: 42px;
  height: 26px;
  border: 1px solid var(--line);
  border-radius: 4px;
  background: #101415;
  color: var(--ink);
  text-align: center;
  font-size: 11px;
}
.histogram {
  display: flex;
  align-items: end;
  gap: 2px;
  height: 72px;
  padding: 8px;
  border: 1px solid var(--line);
  background: #101415;
}
.histogram span {
  flex: 1;
  min-height: 4px;
  background: #6a777c;
}
.metadata {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
  margin: 0;
  font-size: 12px;
}
.metadata dt {
  color: var(--muted);
}
.metadata dd {
  margin: 0;
  color: #d9dfe1;
  text-align: right;
}
.timeline {
  grid-column: 2;
  grid-row: 3;
  min-width: 0;
  border-top: 1px solid var(--line);
  background: #171b1d;
}
.timeline-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 34px;
  padding: 0 12px;
  border-bottom: 1px solid var(--line);
}
.timeline-title {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
}
.timeline-track {
  position: relative;
  height: 61px;
  margin: 0 14px;
}
.timeline-line {
  position: absolute;
  top: 26px;
  right: 0;
  left: 0;
  height: 3px;
  background: var(--surface-soft);
}
.playhead {
  position: absolute;
  top: 8px;
  bottom: 5px;
  left: 34%;
  width: 1px;
  background: var(--accent);
}
.playhead::before {
  position: absolute;
  top: 0;
  left: -4px;
  width: 9px;
  height: 9px;
  border-radius: 2px;
  background: var(--accent);
  content: "";
}
.frame-marker {
  position: absolute;
  top: 17px;
  width: 20px;
  height: 20px;
  border: 2px solid #8b989d;
  border-radius: 3px;
  background: #252c2f;
  transform: translateX(-50%);
}
.frame-marker.ready {
  border-color: var(--accent);
}
.toast {
  position: fixed;
  right: 18px;
  bottom: 18px;
  z-index: 20;
  max-width: 340px;
  padding: 10px 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: #22292c;
  color: var(--ink);
  box-shadow: 0 8px 28px rgb(0 0 0 / 45%);
  opacity: 0;
  pointer-events: none;
  transform: translateY(8px);
  transition: 160ms ease;
}
.toast.visible {
  opacity: 1;
  transform: translateY(0);
}
.toast.error {
  border-color: var(--danger);
}
.empty-list {
  padding: 12px 4px;
  color: var(--muted);
  font-size: 12px;
  line-height: 18px;
}
@media (max-width: 1100px) {
  :root {
    --sidebar: 220px;
    --inspector: 292px;
  }
  .project-name {
    display: none;
  }
}
@media (max-width: 820px) {
  .app {
    grid-template-columns: 1fr;
    grid-template-rows: 54px auto minmax(360px, 1fr) auto auto;
  }
  .topbar {
    grid-column: 1;
    grid-row: 1;
  }
  .brand {
    width: auto;
  }
  .sidebar {
    grid-column: 1;
    grid-row: 2;
    border-right: 0;
    border-bottom: 1px solid var(--line);
  }
  .sidebar-section:not(:first-child) {
    display: none;
  }
  .upload-zone {
    min-height: 76px;
  }
  .workspace {
    grid-column: 1;
    grid-row: 3;
  }
  .timeline {
    grid-column: 1;
    grid-row: 4;
  }
  .inspector {
    grid-column: 1;
    grid-row: 5;
    border-top: 1px solid var(--line);
    border-left: 0;
  }
  .control-section {
    padding: 14px 16px;
  }
  .top-actions .command-button span {
    display: none;
  }
}
</style>
</head>
<body>
<main class="app">
  <header class="topbar">
    <div class="brand">
      <div class="brand-mark">P</div>
      <div class="brand-copy">
        <h1 class="brand-title">宠物短片调色台</h1>
        <div class="brand-subtitle">Pet Color Workbench</div>
      </div>
    </div>
    <div class="project-name" id="projectName">未选择短片</div>
    <div class="top-actions">
      <button class="icon-button" id="refreshButton" title="刷新数据" aria-label="刷新数据">↻</button>
      <button class="command-button" id="previewButton"><span>刷新预览</span></button>
      <button class="command-button primary" id="saveButton"><span>保存调色</span></button>
    </div>
  </header>
  <aside class="sidebar">
    <section class="sidebar-section">
      <h2 class="section-label">素材</h2>
      <label class="upload-zone" for="fileInput">
        <input id="fileInput" type="file" accept="video/mp4,video/quicktime,video/webm" hidden>
        <span>
          <span class="upload-icon">＋</span>
          <span class="upload-main">选择猫狗视频</span>
          <span class="upload-note">MP4、MOV 或 WebM，最长 3 分钟</span>
        </span>
      </label>
    </section>
    <section class="sidebar-section">
      <h2 class="section-label">素材库</h2>
      <input class="search" id="searchInput" type="search" placeholder="搜索文件名">
      <div class="filter-row">
        <select class="select" id="petFilter" aria-label="宠物类型">
          <option value="">全部宠物</option>
          <option value="cat">猫</option>
          <option value="dog">狗</option>
        </select>
        <select class="select" id="stateFilter" aria-label="素材状态">
          <option value="">全部状态</option>
          <option value="ready">可用</option>
          <option value="expired">已过期</option>
        </select>
      </div>
      <div class="clip-list" id="clipList"></div>
    </section>
  </aside>
  <section class="workspace">
    <div class="workspace-toolbar">
      <div class="segmented" role="tablist" aria-label="预览模式">
        <button class="segment active" data-mode="after">调色后</button>
        <button class="segment" data-mode="split">对比</button>
        <button class="segment" data-mode="before">原片</button>
      </div>
      <button class="icon-button" id="fitButton" title="适应画布" aria-label="适应画布">⌗</button>
      <span class="zoom-copy">适应窗口</span>
    </div>
    <div class="stage">
      <div class="viewer" id="viewer">
        <img class="viewer-photo" id="previewImage" alt="短片预览帧">
        <div class="viewer-empty" id="viewerEmpty">
          <span><strong>等待素材</strong>上传或选择短片后开始调色</span>
        </div>
        <div class="viewer-overlay"><span class="status-dot"></span><span id="viewerStatus">工作台就绪</span></div>
      </div>
    </div>
  </section>
  <aside class="inspector">
    <div class="inspector-header">色彩调整</div>
    <section class="control-section">
      <h2 class="section-label">场景预设</h2>
      <div class="preset-grid" id="presetGrid">
        <button class="preset-button active" data-preset="indoor">▤<strong>室内</strong></button>
        <button class="preset-button" data-preset="outdoor">☀<strong>户外</strong></button>
        <button class="preset-button" data-preset="night">◐<strong>夜间</strong></button>
      </div>
    </section>
    <section class="control-section">
      <h2 class="section-label">基础调整</h2>
      <div id="basicControls"></div>
    </section>
    <section class="control-section">
      <h2 class="section-label">亮度分布</h2>
      <div class="histogram" id="histogram"></div>
    </section>
    <section class="control-section">
      <h2 class="section-label">素材信息</h2>
      <dl class="metadata" id="metadata"></dl>
    </section>
  </aside>
  <section class="timeline">
    <div class="timeline-toolbar">
      <span class="timeline-title">预览帧</span>
      <span class="zoom-copy" id="timecode">00:00.000</span>
    </div>
    <div class="timeline-track" id="timelineTrack">
      <div class="timeline-line"></div>
      <div class="playhead"></div>
    </div>
  </section>
</main>
<div class="toast" id="toast" role="status"></div>
<script>
const state = {
  clips: [],
  selected: null,
  grade: null,
  frames: [],
  dirty: false,
  mode: "after"
};
const controls = [
  { key: "exposure", label: "曝光", min: -100, max: 100 },
  { key: "saturation", label: "饱和度", min: -100, max: 100 },
  { key: "temperature", label: "色温", min: -100, max: 100 },
  { key: "contrast", label: "对比度", min: -100, max: 100 },
  { key: "highlights", label: "高光", min: -100, max: 100 },
  { key: "shadows", label: "阴影", min: -100, max: 100 }
];
const elements = {
  clipList: document.querySelector("#clipList"),
  projectName: document.querySelector("#projectName"),
  viewerEmpty: document.querySelector("#viewerEmpty"),
  viewerStatus: document.querySelector("#viewerStatus"),
  previewImage: document.querySelector("#previewImage"),
  presetGrid: document.querySelector("#presetGrid"),
  basicControls: document.querySelector("#basicControls"),
  histogram: document.querySelector("#histogram"),
  metadata: document.querySelector("#metadata"),
  timelineTrack: document.querySelector("#timelineTrack"),
  timecode: document.querySelector("#timecode"),
  toast: document.querySelector("#toast"),
  searchInput: document.querySelector("#searchInput"),
  petFilter: document.querySelector("#petFilter"),
  stateFilter: document.querySelector("#stateFilter"),
  fileInput: document.querySelector("#fileInput")
};
function escapeHTML(value) {
  const node = document.createElement("span");
  node.textContent = String(value ?? "");
  return node.innerHTML;
}
function showToast(message, error = false) {
  elements.toast.textContent = message;
  elements.toast.classList.toggle("error", error);
  elements.toast.classList.add("visible");
  window.clearTimeout(showToast.timer);
  showToast.timer = window.setTimeout(() => elements.toast.classList.remove("visible"), 2600);
}
async function api(path, options = {}) {
  const response = await fetch(path, options);
  if (response.status === 204) {
    return null;
  }
  const payload = await response.json();
  if (!response.ok) {
    throw new Error(payload.error || "请求失败");
  }
  return payload;
}
function buildControls() {
  elements.basicControls.innerHTML = controls.map((control) => '
    <div class="slider-control">
      <label for="control-' + control.key + '">' + control.label + '</label>
      <input id="control-' + control.key + '" data-key="' + control.key + '" type="range" min="' + control.min + '" max="' + control.max + '" value="0">
      <input class="slider-value" data-value="' + control.key + '" type="text" value="0" readonly>
    </div>
  ').join("");
  elements.basicControls.querySelectorAll("input[type=range]").forEach((input) => {
    input.addEventListener("input", () => {
      if (!state.grade) return;
      const key = input.dataset.key;
      state.grade[key] = Number(input.value);
      state.dirty = true;
      const output = elements.basicControls.querySelector('[data-value="' + key + '"]');
      output.value = Number(input.value) > 0 ? '+' + input.value : input.value;
      renderHistogram();
      elements.viewerStatus.textContent = "参数已修改";
    });
  });
}
function renderClips() {
  if (!state.clips.length) {
    elements.clipList.innerHTML = '<div class="empty-list">还没有素材。上传猫或狗的短片，工作台会创建默认调色会话。</div>';
    return;
  }
  elements.clipList.innerHTML = state.clips.map((item) => {
    const clip = item.clip;
    const active = state.selected && state.selected.clip.id === clip.id;
    const pet = clip.pet_kind === "cat" ? "猫" : "狗";
    return '
      <button class="clip-item ' + (active ? "active" : "") + '" data-clip="' + escapeHTML(clip.id) + '">
        <span class="clip-thumb">' + (clip.pet_kind === "cat" ? "猫" : "狗") + '</span>
        <span class="clip-meta">
          <span class="clip-name">' + escapeHTML(clip.source_name) + '</span>
          <span class="clip-detail">' + pet + ' · ' + formatDuration(clip.duration_ms) + '</span>
        </span>
      </button>
    ';
  }).join("");
  elements.clipList.querySelectorAll("[data-clip]").forEach((button) => {
    button.addEventListener("click", () => selectClip(button.dataset.clip));
  });
}
function renderSelection() {
  if (!state.selected || !state.grade) {
    elements.projectName.textContent = "未选择短片";
    elements.viewerEmpty.style.display = "grid";
    elements.previewImage.style.display = "none";
    elements.metadata.innerHTML = "";
    return;
  }
  const clip = state.selected.clip;
  elements.projectName.textContent = clip.source_name;
  elements.viewerEmpty.style.display = "none";
  elements.viewerStatus.textContent = '修订 ' + state.grade.revision + ' · ' + formatDuration(clip.duration_ms);
  controls.forEach((control) => {
    const input = elements.basicControls.querySelector('[data-key="' + control.key + '"]');
    const output = elements.basicControls.querySelector('[data-value="' + control.key + '"]');
    input.value = state.grade[control.key] || 0;
    output.value = state.grade[control.key] > 0 ? '+' + state.grade[control.key] : state.grade[control.key] || 0;
  });
  elements.presetGrid.querySelectorAll("[data-preset]").forEach((button) => {
    button.classList.toggle("active", button.dataset.preset === state.grade.preset);
  });
  elements.metadata.innerHTML = '
    <dt>宠物</dt><dd>' + (clip.pet_kind === "cat" ? "猫" : "狗") + '</dd>
    <dt>分辨率</dt><dd>' + clip.width + ' × ' + clip.height + '</dd>
    <dt>时长</dt><dd>' + formatDuration(clip.duration_ms) + '</dd>
    <dt>大小</dt><dd>' + formatBytes(clip.size_bytes) + '</dd>
    <dt>状态</dt><dd>' + clip.state + '</dd>
    <dt>修订</dt><dd>' + state.grade.revision + '</dd>
  ';
  renderHistogram();
  renderTimeline();
}
function renderHistogram() {
  const exposure = state.grade ? state.grade.exposure : 0;
  const contrast = state.grade ? state.grade.contrast : 0;
  const values = [12, 18, 27, 38, 52, 66, 78, 84, 89, 82, 72, 58, 43, 31, 20, 13];
  elements.histogram.innerHTML = values.map((base, index) => {
    const centerDistance = Math.abs(index - 8);
    const height = Math.max(5, Math.min(100, base + exposure * 0.24 + contrast * (4 - centerDistance) * 0.08));
    return '<span style="height:' + height + '%"></span>';
  }).join("");
}
function renderTimeline() {
  elements.timelineTrack.querySelectorAll(".frame-marker").forEach((node) => node.remove());
  if (!state.selected) return;
  const duration = state.selected.clip.duration_ms || 1;
  state.frames.forEach((frame) => {
    const marker = document.createElement("span");
    marker.className = 'frame-marker ' + (frame.status === "ready" ? "ready" : "");
    marker.style.left = Math.min(100, Math.max(0, frame.timestamp_ms / duration * 100)) + '%';
    marker.title = formatDuration(frame.timestamp_ms);
    elements.timelineTrack.appendChild(marker);
  });
}
function formatDuration(milliseconds) {
  const total = Math.max(0, Number(milliseconds) || 0);
  const minutes = Math.floor(total / 60000);
  const seconds = Math.floor(total % 60000 / 1000);
  const millis = Math.floor(total % 1000);
  return String(minutes).padStart(2, "0") + ':' + String(seconds).padStart(2, "0") + '.' + String(millis).padStart(3, "0");
}
function formatBytes(bytes) {
  const value = Number(bytes) || 0;
  if (value < 1024) return value + ' B';
  if (value < 1024 * 1024) return (value / 1024).toFixed(1) + ' KB';
  return (value / 1024 / 1024).toFixed(1) + ' MB';
}
async function loadClips() {
  const params = new URLSearchParams();
  if (elements.searchInput.value) params.set("q", elements.searchInput.value);
  if (elements.petFilter.value) params.set("pet", elements.petFilter.value);
  if (elements.stateFilter.value) params.set("state", elements.stateFilter.value);
  const result = await api('/api/clips?' + params);
  state.clips = result.items;
  renderClips();
}
async function selectClip(identifier) {
  try {
    const summary = await api('/api/clips/' + encodeURIComponent(identifier));
    state.selected = summary;
    state.grade = structuredClone(summary.grade);
    state.frames = summary.latest_preview ? [summary.latest_preview] : [];
    state.dirty = false;
    renderClips();
    renderSelection();
  } catch (error) {
    showToast(error.message, true);
  }
}
async function applyPreset(name) {
  if (!state.selected) {
    showToast("请先选择短片", true);
    return;
  }
  try {
    const grade = await api('/api/clips/' + encodeURIComponent(state.selected.clip.id) + '/preset', {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name })
    });
    state.grade = grade;
    state.selected.grade = grade;
    state.dirty = false;
    renderSelection();
    showToast("已应用场景预设");
  } catch (error) {
    showToast(error.message, true);
  }
}
async function saveGrade() {
  if (!state.selected || !state.grade) {
    showToast("请先选择短片", true);
    return;
  }
  const payload = { revision: state.selected.grade.revision };
  controls.forEach((control) => payload[control.key] = state.grade[control.key]);
  try {
    const grade = await api('/api/clips/' + encodeURIComponent(state.selected.clip.id) + '/grade', {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    state.grade = grade;
    state.selected.grade = grade;
    state.dirty = false;
    renderSelection();
    showToast("调色参数已保存");
  } catch (error) {
    showToast(error.message, true);
    if (error.message.includes("conflict")) await selectClip(state.selected.clip.id);
  }
}
async function refreshPreview() {
  if (!state.selected) {
    showToast("请先选择短片", true);
    return;
  }
  if (state.dirty) await saveGrade();
  const duration = state.selected.clip.duration_ms;
  const timestamps = [duration * 0.15, duration * 0.5, duration * 0.85].map(Math.floor);
  try {
    elements.viewerStatus.textContent = "正在规划预览…";
    const result = await api('/api/clips/' + encodeURIComponent(state.selected.clip.id) + '/previews', {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ timestamps_ms: timestamps, prefix: 'refresh-' + state.grade.revision })
    });
    state.frames = result.frames;
    renderTimeline();
    elements.viewerStatus.textContent = state.frames.length + ' 个预览帧已刷新';
    showToast("预览计划已刷新");
  } catch (error) {
    showToast(error.message, true);
  }
}
async function inspectFile(file) {
  if (!file) return;
  const extension = file.name.split(".").pop().toLowerCase();
  if (!["mp4", "mov", "webm"].includes(extension)) {
    showToast("请选择 MP4、MOV 或 WebM 视频", true);
    return;
  }
  const video = document.createElement("video");
  const objectURL = URL.createObjectURL(file);
  video.preload = "metadata";
  video.src = objectURL;
  video.addEventListener("loadedmetadata", async () => {
    const now = Date.now();
    const identifier = 'clip-' + now;
    const petKind = file.name.toLowerCase().includes("dog") ? "dog" : "cat";
    const mediaType = file.type || (extension === "mov" ? "video/quicktime" : 'video/' + extension);
    const checksum = file.size.toString(16) + file.lastModified.toString(16);
    try {
      const summary = await api("/api/clips", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          id: identifier,
          pet_kind: petKind,
          source_name: file.name,
          media_type: mediaType,
          size_bytes: file.size,
          duration_ms: Math.max(250, Math.floor(video.duration * 1000)),
          width: video.videoWidth,
          height: video.videoHeight,
          checksum,
          ttl_minutes: 120
        })
      });
      await loadClips();
      await selectClip(summary.clip.id);
      elements.previewImage.src = objectURL;
      elements.previewImage.style.display = "block";
      elements.viewerEmpty.style.display = "none";
      showToast("素材已加入调色台");
    } catch (error) {
      URL.revokeObjectURL(objectURL);
      showToast(error.message, true);
    }
  });
  video.addEventListener("error", () => {
    URL.revokeObjectURL(objectURL);
    showToast("无法读取视频元数据", true);
  });
}
elements.presetGrid.addEventListener("click", (event) => {
  const button = event.target.closest("[data-preset]");
  if (button) applyPreset(button.dataset.preset);
});
document.querySelectorAll("[data-mode]").forEach((button) => {
  button.addEventListener("click", () => {
    state.mode = button.dataset.mode;
    document.querySelectorAll("[data-mode]").forEach((item) => item.classList.toggle("active", item === button));
    elements.viewerStatus.textContent = button.textContent.trim();
  });
});
document.querySelector("#saveButton").addEventListener("click", saveGrade);
document.querySelector("#previewButton").addEventListener("click", refreshPreview);
document.querySelector("#refreshButton").addEventListener("click", () => loadClips().catch((error) => showToast(error.message, true)));
document.querySelector("#fitButton").addEventListener("click", () => showToast("画布已适应窗口"));
elements.searchInput.addEventListener("input", () => {
  window.clearTimeout(elements.searchInput.timer);
  elements.searchInput.timer = window.setTimeout(() => loadClips().catch((error) => showToast(error.message, true)), 180);
});
elements.petFilter.addEventListener("change", () => loadClips().catch((error) => showToast(error.message, true)));
elements.stateFilter.addEventListener("change", () => loadClips().catch((error) => showToast(error.message, true)));
elements.fileInput.addEventListener("change", () => inspectFile(elements.fileInput.files[0]));
elements.timelineTrack.addEventListener("click", (event) => {
  if (!state.selected) return;
  const bounds = elements.timelineTrack.getBoundingClientRect();
  const ratio = Math.min(1, Math.max(0, (event.clientX - bounds.left) / bounds.width));
  elements.timelineTrack.querySelector(".playhead").style.left = ratio * 100 + '%';
  elements.timecode.textContent = formatDuration(state.selected.clip.duration_ms * ratio);
});
buildControls();
renderHistogram();
loadClips().catch((error) => showToast(error.message, true));
</script>
</body>
</html>`
