# 问题回答
 ## 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
 ### 从工程化角度来说,sql.ErrNoRows是go database 标准库错误，一般建议通过wrap 包装后，返回给调用方，wrap的时候，可以携带一些上下文提供信息，比如某些sql语句信息、调试信息等，方便定位问题。
 ### 从业务角度来说，ErrNoRows不算是一种异常，只是查询不到数据，可以根据业务实际情况来说处理，比如将errNoRows 转换成nil，返回给调用方，或者直接抛给调用方，由调用方决定怎么处理异常。
### 伪代码
``` go
func findUserNameById(id string)  (name string, err error){
    var name string
    err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)
    if err != nil {
        if err == sql.ErrNoRows {
            // 直接返回nil
            return nil,nil
            //或者返回errNoRows
            //return nil,error
        } else {
            return nil,errors.wrap(err,"find user sql exception")
        }
    }
    return name,nil
}
```
