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
    return str.replace(/&amp;#34;/g, '"')
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

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById("date_start").addEventListener('change', get_online_charts);
    document.getElementById("date_end").addEventListener('change', get_online_charts);

    let showChronicles = true;
    document.getElementById("toggleChronicles").addEventListener('change', function() {
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
    get_online_charts();
});

async function get_online_charts() {
    const menuDateStart = document.getElementById('date_start').value;
    const menuDateEnd = document.getElementById('date_end').value;

    document.getElementById("online-stat-current-dates").innerHTML = `
        <span class="text-danger">${menuDateStart}</span> - <span class="text-success">${menuDateEnd}</span>
        `;

    await getChronicles()

    await get_online_chart("online-stat-all-weeks", "online_stat_weeks", "players, avg", menuDateStart, menuDateEnd)
    await get_online_chart("online-stat-daytime", "online_stat_daytime", "players, avg", menuDateStart, menuDateEnd)

    const dateTo = new Date()
    const dateFrom = new Date()
    dateTo.setDate(dateTo.getDate() - 1) // without today
    dateFrom.setDate(dateTo.getDate() - 90)

    await get_online_chart("online-stat-month", "online_stat", "players, avg", format_date(dateFrom), format_date(dateTo))
    await get_online_chart("online-stat-month", "online_stat_max", "players, max", format_date(dateFrom), format_date(dateTo))
}

function get_online_chart(targetId, endpoint, label, dateFrom, dateTo) {
    return new Promise((resolve, reject) => {
        const params = new URLSearchParams({dateFrom: dateFrom, dateTo: dateTo});
        fetch(`/api/${endpoint}?${params}`)
            .then(response => response.json())
            .then(data => {
                const maxOnline = Math.max(...Object.values(data));
                const chart = Chart.getChart(targetId);

                const newLabels = Object.keys(data);
                const newData = Object.values(data);

                const chroniclesInRange = endpoint === 'online_stat_weeks'
                    ? getChroniclesForWeeks(newLabels)
                    : getChroniclesInRange(newLabels, dateFrom, dateTo);

                if (chart) {
                    chart.data.labels = newLabels;
                    let existingDataset = chart.data.datasets.find(ds => ds.label === label);
                    if (existingDataset) {
                        existingDataset.data = newData;
                    } else {
                        chart.data.datasets.push({
                            label: label,
                            data: newData,
                            borderWidth: 1,
                        });
                    }
                    chart.update();
                    resolve();
                    return;
                }

                const activePlugins = [
                    {
                        id: 'chroniclesPlugin',
                        afterDraw: function(chart) {
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
                        afterEvent: function(chart, args) {
                            if (!document.getElementById("toggleChronicles").checked || args.event.type !== 'mousemove') return;
                            if (args.event.type === 'mousemove') {
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
                                    ctx.moveTo(centerX - rectWidth/2 + radius, centerY - rectHeight/2);
                                    ctx.lineTo(centerX + rectWidth/2 - radius, centerY - rectHeight/2);
                                    ctx.quadraticCurveTo(centerX + rectWidth/2, centerY - rectHeight/2,
                                        centerX + rectWidth/2, centerY - rectHeight/2 + radius);
                                    ctx.lineTo(centerX + rectWidth/2, centerY + rectHeight/2 - radius);
                                    ctx.quadraticCurveTo(centerX + rectWidth/2, centerY + rectHeight/2,
                                        centerX + rectWidth/2 - radius, centerY + rectHeight/2);
                                    ctx.lineTo(centerX - rectWidth/2 + radius, centerY + rectHeight/2);
                                    ctx.quadraticCurveTo(centerX - rectWidth/2, centerY + rectHeight/2,
                                        centerX - rectWidth/2, centerY + rectHeight/2 - radius);
                                    ctx.lineTo(centerX - rectWidth/2, centerY - rectHeight/2 + radius);
                                    ctx.quadraticCurveTo(centerX - rectWidth/2, centerY - rectHeight/2,
                                        centerX - rectWidth/2 + radius, centerY - rectHeight/2);
                                    ctx.closePath();
                                    ctx.fill();
                                    ctx.stroke();

                                    ctx.fillStyle = 'white';
                                    ctx.textAlign = 'center';
                                    ctx.textBaseline = 'middle';

                                    events.forEach((event, index) => {
                                        const yPos = centerY - rectHeight/2 + padding + lineHeight/2 + index * lineHeight;
                                        ctx.fillText(event, centerX, yPos);
                                    });

                                    ctx.restore();
                                }
                            }
                        }
                    }
                ];

                if (targetId === 'online-stat-all-weeks') {
                    activePlugins.push(seasonsPlugin);
                }

                new Chart(
                    document.getElementById(targetId),
                    {
                        type: 'line',
                        data: {
                            labels: newLabels,
                            datasets: [{
                                label: label,
                                data: newData,
                                borderWidth: 1
                            }]
                        },
                        options: {
                            plugins: {
                                colors: {
                                    forceOverride: true
                                },
                                tooltip: {
                                    callbacks: {
                                        afterBody: function(context) {
                                            const label = context[0].label;
                                            if (chroniclesInRange.find(c => c.date === label)) {
                                                const chronicle = chroniclesInRange.find(c => c.date === label);
                                                return `Events:\n${chronicle.text.split('|').join('\n')}`;
                                            }
                                            return null;
                                        }
                                    }
                                }
                            },
                            scales: {
                                y: {
                                    suggestedMin: 0,
                                    suggestedMax: maxOnline + 10
                                },
                                x: { 
                                    grid: {
                                        display: false 
                                    }
                                }
                            }
                        },
                        plugins: activePlugins
                    }
                );
                resolve();
            })
            .catch(error => reject(error));
    })
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
    const dates = [];
    const firstDay = new Date(year, 0, 1);
    const firstWeekDay = firstDay.getDay() || 7;
    let firstWeekDate = new Date(year, 0, 1 + (8 - firstWeekDay) % 7);

    if (week > 1) {
        firstWeekDate.setDate(firstWeekDate.getDate() + (week - 1) * 7);
    }

    for (let i = 0; i < 7; i++) {
        const date = new Date(firstWeekDate);
        date.setDate(date.getDate() + i);
        dates.push(date);
    }

    return dates;
}

function formatDateForChronicle(date) {
    return date.toISOString().split('T')[0];
}

let chronicles = {};
function getChronicles() {
    const dateFrom = document.getElementById('date_start').value;
    const dateTo = document.getElementById('date_end').value;

    return new Promise((resolve, reject) => {
        const params = new URLSearchParams({dateFrom: dateFrom, dateTo: dateTo});
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

function getChroniclesInRange(labels, dateFrom, dateTo) {
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

function format_date(date) {
    let month = date.getMonth() + 1
    if (month < 10)
        month = "0" + month
    let day = date.getDate()
    if (day < 10)
        day = "0" + day
    return `${date.getFullYear()}-${month}-${day}`
}
