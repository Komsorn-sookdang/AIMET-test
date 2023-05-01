db.events.insert([
  {
    name: "Event 1",
    date: new Date("2017-01-01"),
    start_time: "10:00",
    end_time: "23:00",
  },
  {
    name: "Event 2",
    date: new Date("2017-01-25"),
    start_time: "10:00",
    end_time: "12:30",
  },
  {
    name: "Event 3",
    date: new Date("2017-01-25"),
    start_time: "12:24",
    end_time: "20:00",
  },
  {
    name: "Event 4",
    date: new Date("2017-01-20"),
    start_time: "12:59",
    end_time: "13:00",
  },
]);

db.events.createIndex({ date: 1 });
db.events.createIndex({ start_time: 1 });
