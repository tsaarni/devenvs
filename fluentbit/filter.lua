function filter(tag, timestamp, record)
    if string.find(record["message"], "FOO") then
        record["facility"] = "log audit"
    end
    return 1, timestamp, record
end
