<script setup>
import HeaderLogo from './components/HeaderLogo.vue'
import SearchForm from './components/SearchForm.vue'
import Table from './components/Table.vue'
import TextBox from './components/TextBox.vue'

import { ref } from 'vue'

const API_URL = `http://localhost:3000/search/`
//-- --------------------------------------------------------------------------------------------------------------
const mails = ref([
  {
    From: '',
    To: '',
    Subject: '',
    Date: '',
    BodyMessage: ''}
  ])
//-- --------------------------------------------------------------------------------------------------------------
const searchStatus = ref('empty')
//-- --------------------------------------------------------------------------------------------------------------
const selectedMail = ref(null)
//-- --------------------------------------------------------------------------------------------------------------
async function fetcData(keyword){
  const keywordStr = String(keyword)
  const url = `${API_URL}${keywordStr}`
  try{
    const resp = await fetch(url)
    if(!resp.ok){
      throw new Error('Bad response',{cause:{resp}})
    }
    mails.value = await resp.json()
  } catch{
    console.error(err.cause.res?.status)
    throw err
  }
  
  //console.pr(mails.value)
  
  //mails.value = await (await fetch(url)).json()
  searchStatus.value = 'found'
}
//-- --------------------------------------------------------------------------------------------------------------
function setSelectedMail(i){
  if(selectedMail.value == mails.value[i]){
    selectedMail.value = null;
    return;
  }
  selectedMail.value = mails.value[i]
}
</script>

<!-- ============================================================================================================== -->

<template>
  <header>
    <HeaderLogo msg="You did it!" />
  </header>

  <main>
    <SearchForm @search-word="fetcData"/>
    <div class="sm:grid-cols-2 sm:flex h-[42rem]">
      <Table :mails-table="mails" :search-status="searchStatus" @select-row="setSelectedMail" />
      <TextBox :mail-textbox="selectedMail"/>
    </div>
    
  </main>
</template>
<!-- ============================================================================================================== -->
