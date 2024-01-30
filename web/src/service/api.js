import Service from "./service";

class ApiService extends Service {
  async login(param) {
    let data = param;
    return this.fetch(`/api/login`).post(data).json();
  }

  async getServerInfo() {
    return this.fetch(`/api/server`).get().json();
  }
  async sendBroadcast() {
    return this.fetch(`/api/server/broadcast`).post().json();
  }
  async shutdownServer(param) {
    let data = param;
    return this.fetch(`/api/server/shutdown`).post(data).json();
  }

  async getPlayerList(param) {
    const query = this.generateQuery(param);
    return this.fetch(`/api/player?${query}`).get().json();
  }
  async getPlayer(param) {
    const { playerUid } = param;
    return this.fetch(`/api/player/${playerUid}`).get().json();
  }
  async kickPlayer(param) {
    const { playerUid } = param;
    return this.fetch(`/api/player/${playerUid}/kick`).post().json();
  }
  async banPlayer(param) {
    const { playerUid } = param;
    return this.fetch(`/api/player/${playerUid}/ban`).post().json();
  }

  async getGuildList() {
    return this.fetch(`/api/guild`).get().json();
  }
  async getGuild(param) {
    const { adminPlayerUid } = param;
    return this.fetch(`/api/guild/${adminPlayerUid}`).get().json();
  }
}

export default ApiService;
