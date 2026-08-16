const elderGrid = document.querySelector("#elders");
const elderCount = document.querySelector("#elder-count");
const elderSelect = document.querySelector("#elder-id");
const form = document.querySelector("#submission-form");
const formStatus = document.querySelector("#form-status");

function escapeHTML(value) {
  return String(value).replace(/[&<>"']/g, (character) => ({
    "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;"
  }[character]));
}

function renderElders(elders) {
  elderCount.textContent = elders.length;
  elderSelect.replaceChildren(...elders.map((elder) => new Option(elder.name, elder.id)));
  elderGrid.innerHTML = elders.map((elder) => `
    <article class="elder-card">
      <div class="portrait">${escapeHTML(elder.name.slice(0, 1))}</div>
      <div class="elder-card-body">
        <div class="card-heading"><h3>${escapeHTML(elder.name)}</h3><span>公开</span></div>
        <p class="profile">${escapeHTML(elder.profile)}</p>
        <div class="memory-list">
          ${elder.stories.map((story) => `<div class="memory"><span>故事</span><strong>${escapeHTML(story.title)}</strong><p>${escapeHTML(story.body)}</p></div>`).join("")}
          ${elder.importantYears.map((item) => `<div class="memory"><span>年份</span><strong>${item.year} · ${escapeHTML(item.label)}</strong></div>`).join("")}
          ${elder.familyMessages.map((item) => `<div class="memory"><span>寄语</span><strong>${escapeHTML(item.author)}</strong><p>${escapeHTML(item.message)}</p></div>`).join("")}
        </div>
      </div>
    </article>`).join("");
}

async function loadElders() {
  const response = await fetch("/api/elders");
  if (!response.ok) throw new Error("无法加载长者资料");
  renderElders(await response.json());
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  formStatus.textContent = "正在提交...";
  const item = {
    externalId: "web-submission",
    elderId: elderSelect.value,
    kind: document.querySelector("#kind").value,
    title: document.querySelector("#title").value,
    content: document.querySelector("#content").value,
    author: document.querySelector("#author").value,
    visibility: "family"
  };
  const response = await fetch("/api/import-batches", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify([item]) });
  formStatus.textContent = response.ok ? "已提交，等待工作人员审核。" : "提交失败，请稍后重试。";
  if (response.ok) form.reset();
});

loadElders().catch((error) => { elderGrid.innerHTML = `<p class="error">${escapeHTML(error.message)}</p>`; });
