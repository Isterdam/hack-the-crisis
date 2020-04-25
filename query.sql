SELECT coalesce(company_id, crossp.id) as comp_id, coalesce(day, crossp.d) as dow, coalesce( sum(booked) / sum(max) ::float, 0) as avg
FROM 
(
    SELECT company_id, date_part('dow', start_time) as day, max, booked from slots s
    LEFT JOIN bookings b ON s.id=slot_id
    WHERE company_id IN (71,72) AND date_part('week', start_time) = (16)
    GROUP BY company_id, start_time, max, booked
) t 
RIGHT JOIN (
    SELECT a.d, c.id 
    FROM ( VALUES (1), (2), (3), (4), (5), (6), (0)) a (d)
    CROSS JOIN (VALUES (71), (72)) c (id)
) crossp ON crossp.d=t.day AND crossp.id=t.company_id
GROUP BY comp_id, dow
ORDER BY comp_id;

SELECT a.d, c.id 
FROM ( VALUES (0), (1), (2), (3), (4), (5), (6)) a (d)
CROSS JOIN (VALUES (71), (72)) c (id);

SELECT coalesce(t.company_id, crossp.id) as comp_id, coalesce(t.dow, crossp.d) as dow, coalesce(t.count, 0) as count
FROM 
(
    SELECT company_id, date_part('dow', start_time) as dow, count(id) as count from slots s
    WHERE company_id IN (71,72) 
    AND booked < max
    AND date_part('week', start_time) = 16
    GROUP BY company_id, dow
) t
RIGHT JOIN (
    SELECT a.d, c.id 
	FROM company c
	CROSS JOIN ( VALUES (1), (2), (3), (4), (5), (6), (0)) a (d)
	WHERE c.id IN (71,72)
) crossp ON crossp.d=t.dow AND crossp.id=t.company_id
ORDER BY comp_id, dow;



SELECT coalesce(company_id, crossp.id::int) as comp_id, coalesce(day, crossp.d) as dow, coalesce( sum(booked) / sum(max) ::float, 0) as avg
	FROM 
	(
		SELECT company_id, date_part('dow', start_time) as day, max, booked from slots s
		LEFT JOIN bookings b ON s.id=slot_id
		WHERE company_id IN (71,72) AND date_part('week', start_time) = (16)
		GROUP BY company_id, start_time, max, booked
	) t 
	RIGHT JOIN (
		SELECT a.d, c.id 
		FROM company c
		CROSS JOIN ( VALUES (1), (2), (3), (4), (5), (6), (0)) a (d)
		WHERE c.id IN (71,72)
	) crossp ON crossp.d=t.day AND crossp.id::int=t.company_id::int
	GROUP BY comp_id, dow;
	ORDER BY comp_id
 coalesce( sum(booked) / sum(max) ::float, 0) as avg





SELECT t.day::date, coalesce(sum(booked) / sum(max)::float, 0) FROM generate_series(current_timestamp :: timestamp, '2020-05-30'::timestamp, '1 day') t(day)
LEFT JOIN slots s ON s.start_time::date=t.day::date AND s.company_id=71
GROUP BY t.day::date

select * from generate_series(current_timestamp :: timestamp, (current_timestamp::date + 7)::timestamp, '1 day');


SELECT coalesce(s.company_id, 71), t.day::date, coalesce(sum(booked) / sum(max)::float, 0) 
FROM generate_series(current_timestamp, current_timestamp + '7 day', '1 day') t(day)
LEFT JOIN slots s ON s.start_time::date=t.day::date AND s.company_id=71
GROUP BY s.company_id, t.day::date;

SELECT t.day::date, s.id  as avg
FROM generate_series(current_timestamp, current_timestamp + '30 day', '1 day') t(day)
LEFT JOIN slots s ON s.start_time::date=t.day::date AND s.company_id=71 AND s.booked < s.max
GROUP BY s.company_id, t.day::date
ORDER BY t.day::date