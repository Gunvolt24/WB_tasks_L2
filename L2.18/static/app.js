const addBtn = document.getElementById('addEventBtn');
const loadBtn = document.getElementById('loadEventsBtn');
const eventsContainer = document.getElementById('eventsContainer');

addBtn.addEventListener('click', async () => {
  const userId = document.getElementById('userId').value;
  const title = document.getElementById('title').value;
  const description = document.getElementById('description').value;
  const startTime = new Date(document.getElementById('startTime').value).toISOString();
  const endTime = new Date(document.getElementById('endTime').value).toISOString();

  const res = await fetch('/create_event', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: userId, title, description, start_time: startTime, end_time: endTime })
  });

  const data = await res.json();
  if (res.ok) {
    alert('Событие добавлено: ' + data.result);
  } else {
    alert('Ошибка: ' + data.error);
  }
});

loadBtn.addEventListener('click', async () => {
  const userId = document.getElementById('viewUserId').value;
  const date = document.getElementById('viewDate').value;
  const period = document.getElementById('period').value;

  let url = '';
  if (period === 'day') url = `/events_for_day?user_id=${userId}&date=${date}`;
  else if (period === 'week') url = `/events_for_week?user_id=${userId}&date=${date}`;
  else if (period === 'month') url = `/events_for_month?user_id=${userId}&date=${date}`;

  const res = await fetch(url);
  const data = await res.json();

  eventsContainer.innerHTML = '';
  if (res.ok && data.events) {
    data.events.forEach(e => {
      const div = document.createElement('div');
      div.className = 'event';
      div.innerHTML = `<strong>${e.title}</strong> (${e.start_time} → ${e.end_time})<br>${e.description}`;
      eventsContainer.appendChild(div);
    });
  } else {
    eventsContainer.innerHTML = 'Ошибка: ' + (data.error || 'не удалось загрузить события');
  }
});
