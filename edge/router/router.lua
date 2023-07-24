local httpl = require("resty.http")
local cjson = require("cjson")

local router = {}

router.fetch_streams = function(stream_name)
  local httpc = httpl.new()
  local url = "http://api:9090/streams/" .. stream_name
  local res, err = httpc:request_uri(url, {
    method = "GET",
  })
  if not res then
    ngx.log(ngx.ERR, "request failed: ", err)
    return nil, "error"
  end
  ngx.log(ngx.INFO, "got status: ", res.status)
  if res.status ~= 200 then
    return nil, "not found"
  end
  local stream = cjson.decode(res.body)
  ngx.log(ngx.INFO, "host: ", stream.host)

  return stream.host .. "/" .. stream.name .. "/" .. stream.manifest, nil
end

return router
