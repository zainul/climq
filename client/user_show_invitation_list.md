
## user_show_invitation_list
- FieldsRequest : int user_id
- FieldsResponse : errors[],[{user_id,invited,confirm,code,invited_status,invitation_link}]

### Sample : 
- Request : user_show_invitation_list,190119141011331201
- Result: [null,[{"user_id":190119141011331201,"invited":"2019-01-19T07:25:04.172664Z","confirm":null,"code":"190119141011331201_0813332412990_20190119","invited_status":1,"invitation_link":"/qr/invitation/Wy0bkKM3FBfqzjURlahXWAHn8QTFUOGQNf-N5xxhI282N-W2e5x1FiP1DwkbkgyrFo4=.png"}]]