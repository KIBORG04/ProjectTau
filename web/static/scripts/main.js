/**
 * Сбрасывает масштаб графика "Онлайн по неделям".
 */
function resetWeeksZoom() {
    const chart = Chart.getChart('online-stat-all-weeks');
    if (chart) chart.resetZoom();
}

/**
 * Определяет время года по номеру недели.
 * @param {number} week - Номер недели (от 1 до 52/53).
 * @returns {string} Название сезона: 'winter', 'spring', 'summer', 'autumn'.
 */
function getSeason(week) {
    if (week >= 10 && week <= 22) return 'spring';
    if (week >= 23 && week <= 35) return 'summer';
    if (week >= 36 && week <= 48) return 'autumn';
    return 'winter';
}

const seasonsPlugin = {
    id: 'seasonsPlugin',
    beforeDraw: (chart) => {
        const { ctx, chartArea: { top, bottom, right }, scales: { x } } = chart;

        if (!chart.data.labels || chart.data.labels.length === 0) {
            return;
        }

        ctx.save();

        const seasonColors = {
            winter: 'rgba(173, 216, 230, 0.15)',
            spring: 'rgba(255, 182, 193, 0.15)',
            summer: 'rgba(144, 238, 144, 0.15)',
            autumn: 'rgba(255, 165, 0, 0.1)'
        };

        let lastYear = null;

        chart.data.labels.forEach((label, index) => {
            const [year, week] = label.split('-').map(Number);

            const currentSeason = getSeason(week);
            ctx.fillStyle = seasonColors[currentSeason];

            const xStart = x.getPixelForValue(index);
            const xEnd = (index < chart.data.labels.length - 1) ? x.getPixelForValue(index + 1) : right;
            const width = xEnd - xStart;

            ctx.fillRect(xStart, top, width, bottom - top);

            if (lastYear !== null && year > lastYear) {
                const lineColor = 'rgba(70, 130, 180, 0.7)';

                ctx.beginPath();
                ctx.moveTo(xStart, top);
                ctx.lineTo(xStart, bottom);
                ctx.lineWidth = 1.2;
                ctx.strokeStyle = lineColor;
                ctx.stroke();

                ctx.save();
                ctx.font = '12px Arial';
                ctx.fillStyle = lineColor;
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';

                ctx.translate(xStart - 7, top + 50);
                ctx.rotate(-Math.PI / 2);
                ctx.fillText('Новый год', 0, 0);

                ctx.restore();
            }
            lastYear = year;
        });

        ctx.restore();
    }
};

function unescape(str) {
    if (!str) return str;
    return str.replace(/&#34;/g, '"')
        .replace(/&amp;/g, '&')
        .replace(/&lt;/g, '<')
        .replace(/&gt;/g, '>')
        .replace(/&quot;/g, '"')
        .replace(/&#039;/g, "'")
}

function get_announce() {
    fetch("/api/random_announce")
        .then(response => response.json())
        .then(data => {
            const author = data.Author ? `<figcaption class="blockquote-footer">${data.Author}</figcaption>` : ""
            const button = `<button id="update_announce" class="btn btn-outline-secondary btn-sm fa-solid fa-arrow-rotate-right float-md-end ms-2"></button>`
            document.getElementById("random_announce").innerHTML = `
             <figure class="text-end">
                        <p class="h5">${button}${data.Title}</p>
                        <p id="announce_contents" class="fs-6"></p>
                         ${author}
             </figure>`;
            document.getElementById("announce_contents").textContent = unescape(data.Content);
            document.getElementById("update_announce").addEventListener('click', get_announce);
        });
}

function get_achievement() {
    fetch("/api/random_achievement")
        .then(response => response.json())
        .then(data => {
            const button = `<button id="update_achievement" class="btn btn-outline-secondary btn-sm fa-solid fa-arrow-rotate-right float-md-end ms-2"></button>`
            document.getElementById("random_achievement").innerHTML = `
             <figure class="text-end">
                        <p class="h5">${button}${data.Title}</p>
                        <p id="achievement_contents" class="fs-6"></p>
                        <figcaption class="blockquote-footer">
                            ${data.Key} as ${data.Name}
                        </figcaption>
             </figure>`;
            document.getElementById("achievement_contents").textContent = unescape(data.Desc);
            document.getElementById("update_achievement").addEventListener('click', get_achievement);
        });
}

function get_flavor() {
    fetch("/api/random_flavor")
        .then(response => response.json())
        .then(data => {
            const button = `<button id="update_flavor" class="btn btn-outline-secondary btn-sm fa-solid fa-arrow-rotate-right float-md-end ms-2"></button>`
            document.getElementById("random_flavor").innerHTML = `
             <figure class="text-end">
                        <p class="h5">${button}${data.Name}</p>
                        <p id="flavor_contents" class="fs-6"></p>
                        <figcaption class="blockquote-footer" style="text-transform:capitalize;">
                            ${data.Gender}, ${data.Age}, ${data.Species}
                        </figcaption>
             </figure>`;
            document.getElementById("flavor_contents").textContent = unescape(data.Flavor);
            document.getElementById("update_flavor").addEventListener('click', get_flavor);
        });
}

function get_last_phrase() {
    fetch("/api/random_last_phrase")
        .then(response => response.json())
        .then(data => {
            const button = `<button id="update_last_phrase" class="btn btn-outline-secondary btn-sm fa-solid fa-arrow-rotate-right float-md-end ms-2"></button>`
            document.getElementById("random_last_phrase").innerHTML = `
             <figure class="text-end">
                        <p class="h5">${button}Перед смертью в ${data.TimeOfDeath}</p>
                        <p class="fs-6"><span style='font-weight: 500;'>${data.Name}</span>: "${data.Phrase}"</p>
                        <figcaption class="blockquote-footer">
                            Round #${data.RoundID}
                        </figcaption>
             </figure>`;
            document.getElementById("update_last_phrase").addEventListener('click', get_last_phrase);
        });
}

// ---- Global state ----
let onlineStatsData = null;
let chronicles = {};

document.addEventListener('DOMContentLoaded', () => {
    // Date inputs still trigger chronicle reload + chart rebuild for date-dependent views
    document.getElementById("date_start").addEventListener('change', onDatesChanged);
    document.getElementById("date_end").addEventListener('change', onDatesChanged);

    let showChronicles = true;
    document.getElementById("toggleChronicles").addEventListener('change', function () {
        showChronicles = this.checked;

        ['online-stat-all-weeks', 'online-stat-daytime', 'online-stat-month'].forEach(id => {
            const chart = Chart.getChart(id);
            if (chart) {
                chart.update();
            }
        });
    });

    get_announce();
    get_achievement();
    get_last_phrase();
    get_flavor();

    loadOnlineStats();
});

/**
 * Loads the pre-calculated JSON and builds all charts.
 */
async function loadOnlineStats() {
    try {
        const response = await fetch('/web/static/data/online_stats.json');
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        onlineStatsData = await response.json();
    } catch (err) {
        console.error('Failed to load online stats JSON:', err);
        return;
    }

    await getChronicles();
    buildAllCharts();
}

/**
 * Called when the user changes dates in the menu bar.
 * Only chronicles need reloading since chart data is from the JSON.
 */
async function onDatesChanged() {
    const menuDateStart = document.getElementById('date_start').value;
    const menuDateEnd = document.getElementById('date_end').value;

    document.getElementById("online-stat-current-dates").innerHTML = `
        <span class="text-danger">${menuDateStart}</span> - <span class="text-success">${menuDateEnd}</span>
        `;

    await getChronicles();
    rebuildChartsAfterChronicleUpdate();
}

/**
 * Build all three chart canvases from the loaded JSON data.
 */
function buildAllCharts() {
    if (!onlineStatsData) return;

    const menuDateStart = document.getElementById('date_start').value;
    const menuDateEnd = document.getElementById('date_end').value;

    document.getElementById("online-stat-current-dates").innerHTML = `
        <span class="text-danger">${menuDateStart}</span> - <span class="text-success">${menuDateEnd}</span>
        `;

    // ---- Chart 1: Online by Weeks (with zoom/pan) ----
    buildWeeksChart();

    // ---- Chart 2: Average online by hours ----
    buildDaytimeChart();

    // ---- Chart 3: Last 90 days ----
    buildLast90DaysChart();
}

/**
 * Rebuild charts that depend on chronicles (i.e. after date change).
 */
function rebuildChartsAfterChronicleUpdate() {
    // Destroy and recreate charts with updated chronicles
    ['online-stat-all-weeks', 'online-stat-daytime', 'online-stat-month'].forEach(id => {
        const chart = Chart.getChart(id);
        if (chart) chart.destroy();
    });
    buildAllCharts();
}

// ======================================================================
//  Chart 1: Online by Weeks (with zoom/pan)
// ======================================================================

/**
 * Calculates Simple Moving Average (SMA).
 * For early elements (index < windowSize), it averages available elements.
 */
function calculateSMA(data, windowSize) {
    if (!data || data.length === 0) return [];
    const result = [];
    let sum = 0;
    for (let i = 0; i < data.length; i++) {
        sum += data[i];
        if (i >= windowSize) {
            sum -= data[i - windowSize];
        }
        const count = Math.min(i + 1, windowSize);
        result.push(sum / count);
    }
    return result;
}

function buildWeeksChart() {
    const weeksData = onlineStatsData.weeks;
    const accuLabels = Object.keys(weeksData.accu).sort();
    const pccuLabels = Object.keys(weeksData.pccu).sort();

    // Merge labels and sort
    const allLabelsSet = new Set([...accuLabels, ...pccuLabels]);
    const labels = Array.from(allLabelsSet).sort((a, b) => {
        const [ay, aw] = a.split('-').map(Number);
        const [by, bw] = b.split('-').map(Number);
        return ay !== by ? ay - by : aw - bw;
    });

    const accuDataRaw = labels.map(l => weeksData.accu[l] || 0);
    const pccuDataRaw = labels.map(l => weeksData.pccu[l] || 0);

    const accuData = calculateSMA(accuDataRaw, 4);
    const pccuData = calculateSMA(pccuDataRaw, 4);

    const maxOnline = Math.max(...accuData, ...pccuData);

    const chroniclesInRange = getChroniclesForWeeks(labels);

    const canvas = document.getElementById('online-stat-all-weeks');

    new Chart(canvas, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'ACCU (avg)',
                    data: accuData,
                    borderWidth: 2,
                    borderColor: 'rgba(54, 162, 235, 1)',
                    backgroundColor: 'rgba(54, 162, 235, 0.1)',
                    fill: true,
                    tension: 0.3,
                    pointRadius: 0,
                    hitRadius: 10,
                    hoverRadius: 4,
                },
                {
                    label: 'PCCU (max)',
                    data: pccuData,
                    borderWidth: 2,
                    borderColor: 'rgba(255, 99, 132, 1)',
                    backgroundColor: 'rgba(255, 99, 132, 0.1)',
                    fill: true,
                    tension: 0.3,
                    pointRadius: 0,
                    hitRadius: 10,
                    hoverRadius: 4,
                }
            ]
        },
        options: {
            responsive: true,
            interaction: {
                mode: 'index',
                intersect: false,
            },
            plugins: {
                colors: { forceOverride: true },
                tooltip: {
                    callbacks: {
                        afterBody: function (context) {
                            const label = context[0].label;
                            const chronicle = chroniclesInRange.find(c => c.date === label);
                            if (chronicle) {
                                return `Events:\n${chronicle.text.split('|').join('\n')}`;
                            }
                            return null;
                        }
                    }
                },
                zoom: {
                    pan: {
                        enabled: true,
                        mode: 'x',
                    },
                    zoom: {
                        wheel: { enabled: true },
                        pinch: { enabled: true },
                        drag: {
                            enabled: true,
                            modifierKey: 'shift',
                        },
                        mode: 'x',
                    },
                },
            },
            scales: {
                y: {
                    beginAtZero: true
                },
                x: {
                    grid: { display: false }
                }
            }
        },
        plugins: [
            seasonsPlugin,
            createChroniclesPlugin(chroniclesInRange),
        ]
    });
}

// ======================================================================
//  Chart 2: Daytime
// ======================================================================

function buildDaytimeChart() {
    const daytimeData = onlineStatsData.daytime;

    // Build sorted labels (0, 2, 4, ..., 22) and format as "HH:00"
    const rawKeys = Object.keys(daytimeData.accu).map(Number).sort((a, b) => a - b);
    const labels = rawKeys.map(h => `${String(h).padStart(2, '0')}:00`);
    const data = rawKeys.map(h => daytimeData.accu[h] || 0);
    const maxOnline = Math.max(...data);

    const canvas = document.getElementById('online-stat-daytime');

    new Chart(canvas, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'players, avg',
                data: data,
                borderWidth: 1,
            }]
        },
        options: {
            plugins: {
                colors: { forceOverride: true },
            },
            scales: {
                y: {
                    beginAtZero: true
                },
                x: {
                    grid: { display: false }
                }
            }
        }
    });
}

// ======================================================================
//  Chart 3: Last 90 days
// ======================================================================

function buildLast90DaysChart() {
    const last90 = onlineStatsData.last_90_days;

    const accuLabels = Object.keys(last90.accu).sort();
    const pccuLabels = Object.keys(last90.pccu).sort();
    const allLabelsSet = new Set([...accuLabels, ...pccuLabels]);
    const labels = Array.from(allLabelsSet).sort();

    const accuData = labels.map(l => last90.accu[l] || 0);
    const pccuData = labels.map(l => last90.pccu[l] || 0);
    const maxOnline = Math.max(...accuData, ...pccuData);

    const chroniclesInRange = getChroniclesInRange(labels);

    const canvas = document.getElementById('online-stat-month');

    new Chart(canvas, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'ACCU (avg)',
                    data: accuData,
                    borderWidth: 1,
                },
                {
                    label: 'PCCU (max)',
                    data: pccuData,
                    borderWidth: 1,
                }
            ]
        },
        options: {
            plugins: {
                colors: { forceOverride: true },
                tooltip: {
                    callbacks: {
                        afterBody: function (context) {
                            const label = context[0].label;
                            const chronicle = chroniclesInRange.find(c => c.date === label);
                            if (chronicle) {
                                return `Events:\n${chronicle.text.split('|').join('\n')}`;
                            }
                            return null;
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true
                },
                x: {
                    grid: { display: false }
                }
            }
        },
        plugins: [createChroniclesPlugin(chroniclesInRange)]
    });
}

// ======================================================================
//  Chronicles (events/important dates)
// ======================================================================

function getChronicles() {
    const dateFrom = document.getElementById('date_start').value;
    const dateTo = document.getElementById('date_end').value;

    return new Promise((resolve, reject) => {
        const params = new URLSearchParams({ dateFrom: dateFrom, dateTo: dateTo });
        fetch(`/api/chronicles_daytime?${params}`)
            .then(response => {
                if (!response.ok) throw new Error("Network response was not ok");
                return response.json();
            })
            .then(data => {
                chronicles = {};

                const eventsByDate = {};

                for (const [key, value] of Object.entries(data)) {
                    const dateOnly = key.split('T')[0];

                    if (!eventsByDate[dateOnly]) {
                        eventsByDate[dateOnly] = [];
                    }

                    eventsByDate[dateOnly].push(value);
                }

                for (const [date, events] of Object.entries(eventsByDate)) {
                    if (events.length === 1) {
                        chronicles[date] = events[0];
                    } else {
                        chronicles[date] = events.join(' | ');
                    }
                }

                resolve();
            })
            .catch(error => {
                reject(error);
            });
    });
}

function getChroniclesForWeeks(weekLabels) {
    const result = [];

    weekLabels.forEach(weekLabel => {
        const [year, week] = weekLabel.split('-').map(Number);
        const weekDates = getDatesOfWeek(year, week);
        const weekEvents = [];

        weekDates.forEach(date => {
            const dateStr = formatDateForChronicle(date);
            if (chronicles[dateStr]) {
                weekEvents.push(chronicles[dateStr]);
            }
        });

        if (weekEvents.length > 0) {
            result.push({
                date: weekLabel,
                text: weekEvents.join('|')
            });
        }
    });

    return result;
}

function getDatesOfWeek(year, week) {
    // По стандарту ISO 8601, 4 января всегда принадлежит первой неделе года
    const jan4 = new Date(year, 0, 4);

    // Получаем день недели для 4 января (1 - понедельник, ..., 7 - воскресенье)
    const dayOfJan4 = jan4.getDay() || 7;

    // Находим дату понедельника первой недели (может уйти на конец декабря предыдущего года)
    const firstMonday = new Date(year, 0, 4 - dayOfJan4 + 1);

    // Сдвигаем дату на нужное количество недель вперед
    const startOfWeek = new Date(firstMonday);
    startOfWeek.setDate(firstMonday.getDate() + (week - 1) * 7);

    // Собираем 7 дней этой недели
    const dates = [];
    for (let i = 0; i < 7; i++) {
        const date = new Date(startOfWeek);
        date.setDate(date.getDate() + i);
        dates.push(date);
    }

    return dates;
}

function formatDateForChronicle(date) {
    return date.toISOString().split('T')[0];
}

function getChroniclesInRange(labels) {
    const result = [];
    labels.forEach(date => {
        if (chronicles[date]) {
            result.push({
                date: date,
                text: chronicles[date]
            });
        }
    });
    return result;
}

/**
 * Creates a Chart.js plugin that draws chronicle lines and hover tooltips.
 */
function createChroniclesPlugin(chroniclesInRange) {
    return {
        id: 'chroniclesPlugin',
        afterDraw: function (chart) {
            if (!document.getElementById("toggleChronicles").checked) return;
            const ctx = chart.ctx;
            const xAxis = chart.scales.x;
            const yAxis = chart.scales.y;

            chroniclesInRange.forEach(chronicle => {
                const xPos = xAxis.getPixelForValue(chronicle.date);

                ctx.save();
                ctx.beginPath();
                ctx.moveTo(xPos, yAxis.top);
                ctx.lineTo(xPos, yAxis.bottom);
                ctx.lineWidth = 2;
                ctx.strokeStyle = 'rgba(128, 128, 128, 0.3)';
                ctx.stroke();
                ctx.restore();
            });
        },
        afterEvent: function (chart, args) {
            if (!document.getElementById("toggleChronicles").checked || args.event.type !== 'mousemove') return;
            const xAxis = chart.scales.x;
            const yAxis = chart.scales.y;
            const ctx = chart.ctx;
            const mouseX = args.event.x;

            let closestChronicle = null;
            let minDistance = Infinity;

            chroniclesInRange.forEach(chronicle => {
                const xPos = xAxis.getPixelForValue(chronicle.date);
                const distance = Math.abs(mouseX - xPos);

                if (distance < 50 && distance < minDistance) {
                    minDistance = distance;
                    closestChronicle = chronicle;
                }
            });

            if (closestChronicle) {
                const xPos = xAxis.getPixelForValue(closestChronicle.date);

                chart.draw();

                ctx.save();
                ctx.beginPath();
                ctx.moveTo(xPos, yAxis.top);
                ctx.lineTo(xPos, yAxis.bottom);
                ctx.lineWidth = 3;
                ctx.strokeStyle = 'rgba(255, 0, 0, 0.7)';
                ctx.stroke();

                const events = closestChronicle.text.split('|');

                ctx.font = '12px Arial';
                const lineHeight = 16;
                let maxWidth = 0;

                events.forEach(event => {
                    const width = ctx.measureText(event).width;
                    maxWidth = Math.max(maxWidth, width);
                });

                const padding = 10;
                const rectWidth = maxWidth + padding * 2;
                const rectHeight = events.length * lineHeight + padding * 2;
                const centerX = chart.width / 2;
                const centerY = 55;

                ctx.fillStyle = 'rgba(0, 0, 0, 0.7)';
                ctx.strokeStyle = 'rgba(255, 255, 255, 0.9)';
                ctx.lineWidth = 1;

                const radius = 5;
                ctx.beginPath();
                ctx.moveTo(centerX - rectWidth / 2 + radius, centerY - rectHeight / 2);
                ctx.lineTo(centerX + rectWidth / 2 - radius, centerY - rectHeight / 2);
                ctx.quadraticCurveTo(centerX + rectWidth / 2, centerY - rectHeight / 2,
                    centerX + rectWidth / 2, centerY - rectHeight / 2 + radius);
                ctx.lineTo(centerX + rectWidth / 2, centerY + rectHeight / 2 - radius);
                ctx.quadraticCurveTo(centerX + rectWidth / 2, centerY + rectHeight / 2,
                    centerX + rectWidth / 2 - radius, centerY + rectHeight / 2);
                ctx.lineTo(centerX - rectWidth / 2 + radius, centerY + rectHeight / 2);
                ctx.quadraticCurveTo(centerX - rectWidth / 2, centerY + rectHeight / 2,
                    centerX - rectWidth / 2, centerY + rectHeight / 2 - radius);
                ctx.lineTo(centerX - rectWidth / 2, centerY - rectHeight / 2 + radius);
                ctx.quadraticCurveTo(centerX - rectWidth / 2, centerY - rectHeight / 2,
                    centerX - rectWidth / 2 + radius, centerY - rectHeight / 2);
                ctx.closePath();
                ctx.fill();
                ctx.stroke();

                ctx.fillStyle = 'white';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';

                events.forEach((event, index) => {
                    const yPos = centerY - rectHeight / 2 + padding + lineHeight / 2 + index * lineHeight;
                    ctx.fillText(event, centerX, yPos);
                });

                ctx.restore();
            }
        }
    };
}
