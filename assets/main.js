
document.addEventListener("DOMContentLoaded", function(event) {
  var loginUserId = $('#app').data('login-user-id')

  // https://qiita.com/huigo/items/0efcf50c17b0b1ee27cb#%E9%96%A2%E6%95%B0%E5%9E%8B%E3%82%B3%E3%83%B3%E3%83%9D%E3%83%BC%E3%83%8D%E3%83%B3%E3%83%88
  Vue.component('nl2br', {
    functional: true,
    props: ['tag','text'],
    render: function (createElement, context) {
      return createElement(context.props.tag,
        context.props.text.split("\n").reduce(function(x,y){
          if (!Array.isArray(x)){
            return [x,createElement('br'),y];
          }
          return x.concat([createElement('br'),y]);
        }));
    }
  });

  window.vm = new Vue({
    el: '#app',
    data: {
      body: "",
      sending: false,
      currentRoomId: loginUserId,
      currentTab: [],
      loginUserId: loginUserId,
    },
    methods: {
      sendChat: function() {
        this.sending = true;
        axios.post('/chat/' + this.currentRoomId, {
          body: this.body
        })
        .then(() => {
          // success
          this.sending = false;
          this.body = ""
          this.getChats(this.currentRoomId); // FIXME
        })
        .catch((error) => {
          console.log(error);
          this.sending = false;
        });
      },
      clickTab: function(roomId) {
        this.currentRoomId = roomId
        this.getChats(roomId)
      },
      getChats: function(roomId) {
        axios.get('/chat/' + roomId)
        .then((response) => {
          this.currentTab = response.data
        })
      },
      onMessage: function(payload) {
        var messageUserId = parseInt(payload.user);
        console.log(payload)
        if (this.loginUserId !== this.messageUserId) {
          this.getChats(this.currentRoomId);
        }
      }
    },
    filters: {
      moment: function (value) {
        if (!value) return ''
        return value.slice(5, 16).replace('-', '/').replace('T', ' ')
      }
    },
    mounted: function() {
      this.getChats(this.currentRoomId)

      var url = new URL(location.href);
      var socket = new WebSocket(`${url.protocol.replace("http", "ws")}//${url.host}/socket`);
      socket.onmessage = e => { const payload = JSON.parse(e.data); (payload.type == "CONNECT") ? (console.log("connect")) : this.onMessage(payload) };
      setInterval(() => socket.send(JSON.stringify({type:"KEEPALIVE"})), 40*1000); // https://devcenter.heroku.com/articles/error-codes#h15-idle-connection
      socket.onerror = e => console.log("[ONERROR]", e);
      socket.onclose = e => console.log("[ONCLOSE]", e);
    }
  })

});
