import Service from './service'

class ApiService extends Service {
  async getServerInfo() {
    return this.fetch(`/api/server/info`).get().json()
  }
  async getPlayerList(param) {
    const query = this.generateQuery(param)
    return this.fetch(`/api/player?${query}`).get().json()
  }
  async handlePlayer(param) {
    const { actionId, type } = param
    return this.fetch(`/player/${actionId}/${type}`).post().json()
  }
  async sendBroadcast() {
    return this.fetch(`/broadcast`).post().json()
  }
}

export default ApiService
