-- noinspection SqlResolve
SELECT
    cid,
    name,
    type,
    [notnull],
    dflt_value,
    pk
FROM pragma_table_info(?)
