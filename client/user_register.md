
## user_register
- FieldsRequest : string fullname,string email,string password,string account_type,string shop_name,string username,string phone,string dob,int    school_id,string invitation_code
- FieldsResponse : errors[],user_id

### Sample : 
- Request : user_register,Merchant1,merchant@gmail.com,1234567890,MERCHANT,toko abadi,tokoabadi,083744744,1993-11-11,1,1_system_20190126085720
- Result: [null,190126085818222419]|||||||