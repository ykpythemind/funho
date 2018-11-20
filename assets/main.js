
document.addEventListener("DOMContentLoaded", function(event) {
  var loginUserId = $('#app').data('login-user-id')

  window.vm = new Vue({
    el: '#app',
    data: {
      body: "",
      sending: false,
      roomId: loginUserId,
      currentTab: [],
    },
    methods: {
      sendChat: function() {
        this.sending = true;
        axios.post('/chat/' + this.roomId, {
          body: this.body
        })
        .then(() => {
          // success
          console.log("success")
          this.sending = false;
          this.body = ""
          this.getChats(this.roomId); // FIXME
        })
        .catch((error) => {
          console.log(error);
          this.sending = false;
        });
      },
      clickTab: function(roomId) {
        this.roomId = roomId
        this.getChats(roomId)
      },
      getChats: function(roomId) {
        axios.get('/chat/' + roomId)
        .then((response) => {
          console.log(response.data)
          this.currentTab = response.data
        })
      }
    },
    filters: {
      moment: function (value) {
        if (!value) return ''
        return value.slice(5, 16).replace('-', '/').replace('T', ' ')
      }
    },
    mounted: function() {
      this.getChats(this.roomId)
    }
  })

});
