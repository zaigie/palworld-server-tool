import Service from "./service";

class ApiService extends Service {
  async login(param) {
    let data = param;
    return this.fetch(`/api/login`).post(data).json();
  }

  async getServerToolInfo() {
    return this.fetch(`/api/server/tool`).get().json();
  }
  async getServerInfo() {
    return this.fetch(`/api/server`).get().json();
  }
  async getServerMetrics() {
    return this.fetch(`/api/server/metrics`).get().json();
  }
  async sendBroadcast(param) {
    let data = param;
    return this.fetch(`/api/server/broadcast`).post(data).json();
  }
  async shutdownServer(param) {
    let data = param;
    return this.fetch(`/api/server/shutdown`).post(data).json();
  }

  async getPlayerList(param) {
    const query = this.generateQuery(param);
    return this.fetch(`/api/player?${query}`).get().json();
  }
  async getOnlinePlayerList() {
    return this.fetch(`/api/online_player`).get().json();
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
  async unbanPlayer(param) {
    const { playerUid } = param;
    return this.fetch(`/api/player/${playerUid}/unban`).post().json();
  }

  async getGuildList() {
    return this.fetch(`/api/guild`).get().json();
  }
  async getGuild(param) {
    const { adminPlayerUid } = param;
    return this.fetch(`/api/guild/${adminPlayerUid}`).get().json();
  }

  async getWhitelist() {
    return this.fetch(`/api/whitelist`).get().json();
  }

  async addWhitelist(param) {
    let data = param;
    return this.fetch(`/api/whitelist`).post(data).json();
  }

  async removeWhitelist(param) {
    let data = param;
    return this.fetch(`/api/whitelist`).delete(data).json();
  }

  async putWhitelist(param) {
    let data = param;
    return this.fetch(`/api/whitelist`).put(data).json();
  }

  async getRconCommands() {
    return this.fetch(`/api/rcon`).get().json();
  }

  async sendRconCommand(param) {
    let data = param;
    return this.fetch(`/api/rcon/send`).post(data).json();
  }

  async addRconCommand(param) {
    let data = param;
    return this.fetch(`/api/rcon`).post(data).json();
  }

  async putRconCommand(uuid, param) {
    let data = param;
    return this.fetch(`/api/rcon/${uuid}`).put(data).json();
  }

  async removeRconCommand(uuid) {
    return this.fetch(`/api/rcon/${uuid}`).delete().json();
  }

  async getBackupList(param) {
    const query = this.generateQuery(param);
    return this.fetch(`/api/backup?${query}`).get().json();
  }

  async removeBackup(uuid) {
    return this.fetch(`/api/backup/${uuid}`).delete().json();
  }

  async downloadBackup(uuid) {
    return this.fetch(`/api/backup/${uuid}`).get().blob();
  }
}

export default ApiService;
