<script setup>
import { stringifyExpression } from '@vue/compiler-core';
import { ref } from 'vue'

defineProps({
    mailsTable: {
        type: Array,
        required: true,
    },
    searchStatus: String,
})

const selectedRowIndex = ref(0)

function selectedRowClass(rowIndex){
    if(!this.selectedRowIndex) return '';
    if(rowIndex == this.selectedRowIndex ) return 'bg-cyan-100 hover:bg-cyan-100';
    
    return '';
}

function setIndexSelectedRow(i){
    if(this.selectedRowIndex==i+1){
        this.selectedRowIndex = 0;
        return;
    }
    this.selectedRowIndex = i+1;
}

</script>

<template>
<div class="sm:w-full m-4">
<div class="relative overflow-x-auto shadow-md sm:rounded-lg border">
    <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400 table-fixed">
        <thead class="text-white font-black uppercase bg-[#1D75B8] dark:bg-gray-700 dark:text-gray-400">
            <tr>
                <th scope="col" class="px-6 py-4 w-1/12">
                    N.
                </th>
                <th scope="col" class="px-6 py-4">
                    From
                </th>
                <th scope="col" class="px-6 py-4">
                    To
                </th>
                <th scope="col" class="px-6 py-4">
                    Subject
                </th>
            </tr>
        </thead>
        <tbody>
            <tr v-if = "searchStatus!='empty'" v-for="(mail,i ) in mailsTable" @click="$emit('selectRow',i); setIndexSelectedRow(i)" class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 cursor-pointer" :class="selectedRowClass(i+1)" >
                <td class="px-6 py-4">
                    {{i+1}}
                </td>
                <th scope="row" class="px-6 py-4 whitespace-nowrap dark:text-white truncate">
                    {{ mail.From }}
                </th>
                <td class="px-6 py-4 truncate">
                    {{ mail.To }}
                </td>
                <td class="px-6 py-4 truncate">
                    {{ mail.Subject }}
                </td>
            </tr>
            <tr v-if = "searchStatus=='empty'" class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600">
                <td class="px-6 py-4">
                    
                </td>
                <th scope="row" class="px-6 py-4 whitespace-nowrap dark:text-white">
                    
                </th>
                <td class="px-6 py-4">
                    
                </td>
                <td class="px-6 py-4">
                    
                </td>
            </tr>
        </tbody>
    </table>
    <nav class="flex items-center justify-between pt-4" aria-label="Table navigation">
        <span class="text-sm font-normal text-gray-500 dark:text-gray-400">Showing <span class="font-semibold text-gray-900 dark:text-white">1-20</span> of <span class="font-semibold text-gray-900 dark:text-white">1000</span></span>
        <ul class="inline-flex items-center -space-x-px">
            <li>
                <a href="#" class="block px-3 py-2 ml-0 leading-tight text-gray-500 bg-white border border-gray-300 rounded-l-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                    <span class="sr-only">Previous</span>
                    <svg class="w-5 h-5" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>
                </a>
            </li>
            <li>
                <a href="#" class="px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">1</a>
            </li>
            <li>
                <a href="#" class="px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">2</a>
            </li>
            <li>
                <a href="#" aria-current="page" class="z-10 px-3 py-2 leading-tight text-blue-600 border border-blue-300 bg-blue-50 hover:bg-blue-100 hover:text-blue-700 dark:border-gray-700 dark:bg-gray-700 dark:text-white">3</a>
            </li>
            <li>
                <a href="#" class="px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">...</a>
            </li>
            <li>
                <a href="#" class="px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">100</a>
            </li>
            <li>
                <a href="#" class="block px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 rounded-r-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                    <span class="sr-only">Next</span>
                    <svg class="w-5 h-5" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"></path></svg>
                </a>
            </li>
        </ul>
    </nav>
</div>
</div>

</template>