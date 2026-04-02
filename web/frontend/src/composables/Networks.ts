import { reactive, computed } from 'vue'

export function NetworkState() {
  const state = reactive({
    page: 1,
    pageSize: 5,

    stats: [
      { label: '总数:', value: '5' },
      { label: '自定义:', value: '3' },
      { label: '内置:', value: '2', variant: 'builtin' }
    ],

    networks: [
      { name: 'bridge', driver: 'bridge', scope: 'local', subnet: '172.17.0.0/16', gateway: '172.17.0.1', containers: 3, removable: false },
      { name: 'app-network', driver: 'bridge', scope: 'local', subnet: '192.168.1.0/24', gateway: '192.168.1.1', containers: 5, removable: true },
      { name: 'host', driver: 'host', scope: 'local', subnet: '-', gateway: '-', containers: 0, removable: false }
    ]
  })

  const pagedNetworks = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return state.networks.slice(start, start + state.pageSize)
  })

  return {
    state,
    pagedNetworks
  }
}