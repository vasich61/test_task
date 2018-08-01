-- Тестовый проект mail.ru
-- Запрос из БД, помещение элемента в БД Tarantool


require('strict').on()


SPACE_NAME = 'mailru_test'

local console = require('console')
console.listen('127.0.0.1:3312')


box.cfg {
    listen = 4502
    --log = 'tarantool.txt'
}
--local log = require('log')


local s = box.space[SPACE_NAME] -- test space
if not s then
    s = box.schema.space.create(SPACE_NAME)
    box.space[SPACE_NAME]:create_index('primary', {type = 'tree', parts = {1, 'unsigned'}})
    for i = 1, 25 do
        s:auto_increment({5*i})
    end
    box.schema.user.grant('guest','read,write,execute', 'universe')
    print(string.format('База данных %q создана и проинициализирована', SPACE_NAME))
else
    print(string.format('База данных %q уже существует', SPACE_NAME))
end


function set_value(val)
    local res = s:auto_increment({val})
    print(res)
    return res[1]
end


function get_value(id)
    return s:get({id})[2]
end


function list()
    return s:select()
end