function mess() {
    var iclass =$('#iclass').val();
        $.ajax({
            async : false,
            url : 'http://139.199.94.141/goapi/zhaoxin.php',
            type : "GET",
            dataType : 'jsonp',
            jsonp : 'callback',
            data : {
                number : iclass
            },
            timeout : 20000,
            success : function(data) {
                if (data.name  != 'error') {
                    var cont = data.name  + ' - ' + data.sex + ' - ' + data.iclass;
                    $('#name').val(cont);
                } else {
                    $('#name').val('');
                }
            }

        });
}

function test() {

    var name = $('#name').val();//将刚刚提交的name值取出
    var iclass =$('#iclass').val();
    var phone =$('#phone').val();
    var message=$('#message').val();
    if (name == "") {
        layer.open({
            style: 'font-size:30px;color:#FF5722;'
            ,content: '该学号不存在'
        });
    } else if (phone == "" || message == "") {
        layer.open({
            style: 'font-size:30px;color:#FF5722;'
            ,content: '请认真填写信息'
        });
    } else {
        layer.open({
            type: 2
            ,shadeClose: false
        });
        $.post('/login', {
            name: name,
            phone: phone,
            iclass: iclass,
            message: message,
        }, function (data) {
            layer.closeAll();
            if(data=="success"){
                layer.open({
                    style: 'font-size:30px;color:#5FB878;'
                    ,content: '报名成功'
                });
            } else {
                layer.open({
                    style: 'font-size:30px;color:#FF5722;'
                    ,content: '报名失败，请重试'
                });
            }
        });
    }
}