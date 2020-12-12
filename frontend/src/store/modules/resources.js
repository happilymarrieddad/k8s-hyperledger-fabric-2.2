import axios from '@/utils/axios';

const baseUrl = process.env.VUE_APP_apiurl;

export default {
    namespaced: true,
    state: () => ({
        items: [],
    }),
    mutations: {},
    actions: {
        async getAll({ state }) {
            let [res, err] = await axios.get(`${baseUrl}/resources`)
            if (err) {
                return [null,err];
            }

            state.items = res;

            return [res];
        },
        async getOne({}, id = 0) {
            let [res, err] = await axios.get(`${baseUrl}/resources/${id}`)
            if (err) {
                return [null,err];
            }

            return [res];
        },
        async createOne({}, newObj) {
            let [res, err] = await axios.post(`${baseUrl}/resources`, newObj)
            if (err) {
                return [null,err];
            }

            return [res];
        },
        async updateOne({}, existingObj) {
            let [res, err] = await axios.put(`${baseUrl}/resources`, existingObj.id, existingObj)
            if (err) {
                return [null,err];
            }

            return [res];
        },
        async deleteOne({}, id = 0) {
            let [res, err] = await axios.delete(`${baseUrl}/resources`, id)
            if (err) {
                return [null,err];
            }

            return [res];
        },
        async getTransactions({}, id = 0) {
            let [res, err] = await axios.get(`${baseUrl}/resource_types/${id}/transactions`)
            if (err) {
                return [null,err];
            }

            return [res];
        }
    },
    getters: {}
}
