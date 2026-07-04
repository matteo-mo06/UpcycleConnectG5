<template>
    <UserLayout>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1
                    class="text-3xl font-bold text-gray-800"
                    style="font-family: var(--font-family-title)"
                >
                    Forum
                </h1>
                <p class="text-sm text-gray-400 mt-1">
                    Échangez avec la communauté upcycling
                </p>
            </div>
            <button
                v-if="canCreateTopic"
                @click="showNewTopic = true"
                class="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-dark transition-colors"
            >
                <svg
                    class="w-4 h-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.8"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M12 4v16m8-8H4"
                    />
                </svg>
                Nouveau sujet
            </button>
        </div>

        <div class="flex gap-3 mb-6">
            <div class="flex-1 relative">
                <svg
                    class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.8"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M21 21l-4.35-4.35M17 11A6 6 0 1 0 5 11a6 6 0 0 0 12 0z"
                    />
                </svg>
                <input
                    v-model="search"
                    @input="debouncedFetch"
                    type="text"
                    placeholder="Rechercher un sujet…"
                    class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
            </div>
        </div>

        <div v-if="loading" class="text-center py-12 text-gray-400 text-sm">
            Chargement…
        </div>

        <div
            v-else-if="topics.length === 0"
            class="bg-white rounded-2xl shadow-sm p-12 text-center"
        >
            <svg
                class="w-12 h-12 text-gray-200 mx-auto mb-3"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
            </svg>
            <p class="text-gray-400 text-sm">
                Aucun sujet pour le moment. Soyez le premier à lancer une
                discussion !
            </p>
        </div>

        <div v-else class="space-y-3">
            <div
                v-for="topic in topics"
                :key="topic.id"
                class="bg-white rounded-2xl shadow-sm px-5 py-4 flex items-center gap-4 hover:shadow-md transition-shadow cursor-pointer"
                @click="openTopic(topic)"
            >
                <div
                    class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0"
                >
                    <svg
                        class="w-5 h-5 text-primary"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="1.8"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                        />
                    </svg>
                </div>
                <div class="flex-1 min-w-0">
                    <h3 class="font-semibold text-gray-800 text-sm truncate">
                        {{ topic.title }}
                    </h3>
                    <p class="text-xs text-gray-400 mt-0.5">
                        Par {{ topic.author_name }} ·
                        {{ formatDate(topic.created_at) }}
                    </p>
                </div>
                <div
                    class="flex items-center gap-4 flex-shrink-0 text-xs text-gray-400"
                >
                    <span class="flex items-center gap-1">
                        <svg
                            class="w-3.5 h-3.5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                            />
                        </svg>
                        {{ topic.replies_count }}
                    </span>
                </div>
            </div>
        </div>

        <Pagination
            v-if="total > LIMIT"
            :page="page"
            :total="total"
            :limit="LIMIT"
            @update:page="changePage"
        />

        <div
            v-if="showNewTopic"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="showNewTopic = false"
        >
            <div
                class="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden"
            >
                <div
                    class="flex items-center justify-between px-6 py-4 border-b border-gray-100"
                >
                    <h2
                        class="font-semibold text-gray-800"
                        style="font-family: var(--font-family-title)"
                    >
                        Nouveau sujet
                    </h2>
                    <button
                        @click="showNewTopic = false"
                        class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors"
                    >
                        <svg
                            class="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>
                <div class="px-6 py-5 space-y-4">
                    <div>
                        <label
                            class="block text-xs font-medium text-gray-600 mb-1.5"
                            >Titre</label
                        >
                        <input
                            v-model="newTopic.title"
                            type="text"
                            placeholder="Titre du sujet…"
                            class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                        />
                    </div>
                    <div>
                        <label
                            class="block text-xs font-medium text-gray-600 mb-1.5"
                            >Message</label
                        >
                        <textarea
                            v-model="newTopic.body"
                            rows="5"
                            placeholder="Décrivez votre sujet…"
                            class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                        />
                    </div>
                    <p v-if="topicError" class="text-xs text-red-500">
                        {{ topicError }}
                    </p>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex gap-3">
                    <button
                        @click="showNewTopic = false"
                        class="flex-1 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition-colors"
                    >
                        Annuler
                    </button>
                    <button
                        @click="createTopic"
                        :disabled="topicLoading"
                        class="flex-1 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                    >
                        {{ topicLoading ? "Publication…" : "Publier" }}
                    </button>
                </div>
            </div>
        </div>

        <div
            v-if="openedTopic || openedTopicLoading || openedTopicError"
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
            @click.self="closeTopic"
        >
            <div
                class="bg-white rounded-2xl shadow-xl w-full max-w-2xl max-h-[85vh] flex flex-col overflow-hidden"
            >
                <div
                    class="flex items-center gap-3 px-6 py-4 border-b border-gray-100 flex-shrink-0"
                >
                    <div class="flex-1 min-w-0 flex items-center gap-2">
                        <template v-if="editTopic">
                            <input
                                v-model="editTitle"
                                @keydown.enter="saveTopicTitle"
                                @keydown.escape="editTopic = false"
                                class="flex-1 px-2 py-1 rounded-lg border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                            />
                            <button
                                @click="saveTopicTitle"
                                :disabled="!editTitle.trim()"
                                class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium hover:bg-primary-dark transition-colors disabled:opacity-50 flex-shrink-0"
                            >
                                Sauvegarder
                            </button>
                            <button
                                @click="editTopic = false"
                                class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors flex-shrink-0"
                            >
                                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            </button>
                        </template>
                        <template v-else>
                            <h2
                                class="font-semibold text-gray-800 text-sm truncate"
                                style="font-family: var(--font-family-title)"
                            >
                                {{ openedTopic?.title ?? "…" }}
                            </h2>
                            <div v-if="openedTopic" class="flex items-center gap-0.5 flex-shrink-0">
                                <button
                                    v-if="canEditTopic"
                                    @click="startEditTopic"
                                    title="Modifier le titre"
                                    class="p-1 rounded text-gray-300 hover:text-primary hover:bg-gray-100 transition-colors"
                                >
                                    <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                                    </svg>
                                </button>
                                <button
                                    v-if="canDeleteTopic"
                                    @click="deleteTopic"
                                    title="Supprimer le sujet"
                                    class="p-1 rounded text-gray-300 hover:text-red-500 hover:bg-red-50 transition-colors"
                                >
                                    <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                    </svg>
                                </button>
                                <button
                                    v-if="canReportContent(openedTopic.id_author)"
                                    @click="openReport('topic', openedTopic.id, openedTopic.title)"
                                    title="Signaler le sujet"
                                    class="p-1 rounded text-gray-300 hover:text-orange-500 hover:bg-orange-50 transition-colors"
                                >
                                    <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                        <path stroke-linecap="round" stroke-linejoin="round" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9" />
                                    </svg>
                                </button>
                            </div>
                        </template>
                    </div>
                    <button
                        @click="closeTopic"
                        class="p-1.5 rounded-lg text-gray-400 hover:bg-gray-100 transition-colors flex-shrink-0"
                    >
                        <svg
                            class="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
                    <div
                        v-if="openedTopicLoading"
                        class="text-center text-gray-400 text-sm py-6"
                    >
                        Chargement…
                    </div>
                    <div
                        v-else-if="openedTopicError"
                        class="text-center text-red-400 text-sm py-6"
                    >
                        {{ openedTopicError }}
                    </div>

                    <template v-else-if="openedTopic">
                        <div v-if="topicBody" class="bg-gray-50 rounded-xl p-4">
                            <div class="flex items-start justify-between gap-3">
                                <div class="flex-1 min-w-0">
                                    <div class="flex items-center gap-2 mb-2">
                                        <div
                                            class="w-7 h-7 rounded-full bg-primary flex items-center justify-center text-white text-xs font-semibold flex-shrink-0"
                                        >
                                            {{ topicBody.author_name?.[0] ?? "?" }}
                                        </div>
                                        <span class="text-sm font-medium text-gray-800">{{
                                            topicBody.author_name
                                        }}</span>
                                        <span class="text-xs text-gray-400"
                                            >· {{ formatDate(topicBody.created_at) }}</span
                                        >
                                    </div>
                                    <template v-if="editPost === topicBody.id">
                                        <textarea
                                            v-model="editBody"
                                            rows="4"
                                            @keydown.escape="editPost = null"
                                            class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                                        />
                                        <div class="flex gap-2 mt-2">
                                            <button
                                                @click="editPost = null"
                                                class="px-3 py-1.5 rounded-lg border border-gray-200 text-xs text-gray-600 hover:bg-gray-50 transition-colors"
                                            >
                                                Annuler
                                            </button>
                                            <button
                                                @click="savePostBody(topicBody)"
                                                :disabled="!editBody.trim()"
                                                class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                                            >
                                                Sauvegarder
                                            </button>
                                        </div>
                                    </template>
                                    <p
                                        v-else
                                        class="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap"
                                    >
                                        {{ topicBody.body }}
                                    </p>
                                </div>
                                <div
                                    v-if="editPost !== topicBody.id"
                                    class="flex items-center gap-0.5 flex-shrink-0 pt-0.5"
                                >
                                    <button
                                        v-if="canEditPost(topicBody)"
                                        @click="startEditPost(topicBody)"
                                        title="Modifier"
                                        class="p-1 rounded text-gray-300 hover:text-primary hover:bg-gray-100 transition-colors"
                                    >
                                        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                                        </svg>
                                    </button>
                                    <button
                                        v-if="canDeletePost(topicBody)"
                                        @click="deletePost(topicBody)"
                                        title="Supprimer"
                                        class="p-1 rounded text-gray-300 hover:text-red-500 hover:bg-red-50 transition-colors"
                                    >
                                        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        </div>

                        <template v-if="topicReplies.length > 0">
                            <div class="border-t border-gray-100 pt-1">
                                <p class="text-xs text-gray-400 mb-3">
                                    {{ topicReplies.length }} réponse{{
                                        topicReplies.length > 1 ? "s" : ""
                                    }}
                                </p>
                            </div>
                            <div
                                v-for="post in topicReplies"
                                :key="post.id"
                                :class="[
                                    'space-y-1',
                                    post.id_parent_post ? 'ml-6' : '',
                                ]"
                            >
                                <div
                                    v-if="
                                        post.id_parent_post &&
                                        getParentPost(post.id_parent_post)
                                    "
                                    class="flex items-start gap-1.5 text-xs text-gray-400 mb-1"
                                >
                                    <svg
                                        class="w-3 h-3 mt-0.5 flex-shrink-0"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"
                                        />
                                    </svg>
                                    <div class="border-l-2 border-gray-200 pl-2 min-w-0">
                                        <span class="font-medium text-gray-500">{{
                                            getParentPost(post.id_parent_post).author_name
                                        }}</span>
                                        <p class="truncate italic mt-0.5">{{
                                            getParentPost(post.id_parent_post).body?.slice(0, 80)
                                        }}{{
                                            (getParentPost(post.id_parent_post).body?.length ?? 0) > 80 ? "…" : ""
                                        }}</p>
                                    </div>
                                </div>
                                <div
                                    class="flex items-start justify-between gap-3 bg-white rounded-xl border border-gray-100 p-3"
                                >
                                    <div class="flex-1 min-w-0">
                                        <div
                                            class="flex items-center gap-2 mb-1"
                                        >
                                            <div
                                                class="w-6 h-6 rounded-full bg-secondary/20 flex items-center justify-center text-secondary text-xs font-semibold flex-shrink-0"
                                            >
                                                {{
                                                    post.author_name?.[0] ?? "?"
                                                }}
                                            </div>
                                            <span
                                                class="text-xs font-medium text-gray-700"
                                                >{{ post.author_name }}</span
                                            >
                                            <span class="text-xs text-gray-400"
                                                >·
                                                {{
                                                    formatDate(post.created_at)
                                                }}</span
                                            >
                                        </div>
                                        <template v-if="editPost === post.id">
                                            <textarea
                                                v-model="editBody"
                                                rows="3"
                                                @keydown.escape="editPost = null"
                                                class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30 resize-none"
                                            />
                                            <div class="flex gap-2 mt-2">
                                                <button
                                                    @click="editPost = null"
                                                    class="px-3 py-1.5 rounded-lg border border-gray-200 text-xs text-gray-600 hover:bg-gray-50 transition-colors"
                                                >
                                                    Annuler
                                                </button>
                                                <button
                                                    @click="savePostBody(post)"
                                                    :disabled="!editBody.trim()"
                                                    class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                                                >
                                                    Sauvegarder
                                                </button>
                                            </div>
                                        </template>
                                        <p
                                            v-else
                                            class="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap pl-8"
                                        >
                                            {{ post.body }}
                                        </p>
                                    </div>
                                    <div
                                        v-if="editPost !== post.id"
                                        class="flex items-center gap-2 flex-shrink-0 pt-0.5"
                                    >
                                        <button
                                            v-if="canCreatePost"
                                            @click="startReply(post)"
                                            class="text-xs text-gray-400 hover:text-primary transition-colors"
                                        >
                                            Répondre
                                        </button>
                                        <button
                                            v-if="canEditPost(post)"
                                            @click="startEditPost(post)"
                                            title="Modifier"
                                            class="p-1 rounded text-gray-300 hover:text-primary hover:bg-gray-100 transition-colors"
                                        >
                                            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                                            </svg>
                                        </button>
                                        <button
                                            v-if="canDeletePost(post)"
                                            @click="deletePost(post)"
                                            title="Supprimer"
                                            class="p-1 rounded text-gray-300 hover:text-red-500 hover:bg-red-50 transition-colors"
                                        >
                                            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                            </svg>
                                        </button>
                                        <button
                                            v-if="canReportContent(post.id_author)"
                                            @click="openReport('post', post.id, post.author_name)"
                                            title="Signaler"
                                            class="p-1 rounded text-gray-300 hover:text-orange-500 hover:bg-orange-50 transition-colors"
                                        >
                                            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                                <path stroke-linecap="round" stroke-linejoin="round" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </template>
                    </template>
                </div>

                <div
                    v-if="canCreatePost && openedTopic"
                    class="px-6 py-4 border-t border-gray-100 flex-shrink-0"
                >
                    <div
                        v-if="replyingTo"
                        class="flex items-center justify-between text-xs text-gray-500 bg-gray-50 rounded-lg px-3 py-1.5 mb-2"
                    >
                        <span
                            >En réponse à
                            <span class="font-medium text-gray-700">{{
                                replyingTo.author_name
                            }}</span></span
                        >
                        <button
                            @click="cancelReply"
                            class="text-gray-400 hover:text-gray-600 transition-colors ml-2"
                        >
                            <svg
                                class="w-3.5 h-3.5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M6 18L18 6M6 6l12 12"
                                />
                            </svg>
                        </button>
                    </div>
                    <div class="flex gap-3">
                        <input
                            v-model="replyBody"
                            type="text"
                            :placeholder="
                                replyingTo
                                    ? `Répondre à ${replyingTo.author_name}…`
                                    : 'Votre réponse…'
                            "
                            class="flex-1 px-3 py-2 rounded-xl border border-gray-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                            @keydown.enter="submitReply"
                        />
                        <button
                            @click="submitReply"
                            :disabled="replyLoading || !replyBody.trim()"
                            class="px-4 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-dark transition-colors disabled:opacity-50"
                        >
                            Envoyer
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <ReportModal
            v-model="showReport"
            :content-type="reportTarget?.type"
            :content-id="reportTarget?.id"
            :content-title="reportTarget?.title"
        />
    </UserLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue";
import UserLayout from "@/Layouts/UserLayout.vue";
import Pagination from "@/Components/Pagination.vue";
import ReportModal from "@/Components/ReportModal.vue";
import api from "@/api.js";
import { useAuthStore } from "@/stores/auth.js";
import { formatDate } from "@/utils.js";

const auth = useAuthStore();

const LIMIT = 20;

const topics = ref([]);
const loading = ref(false);
const search = ref("");
const page = ref(1);
const total = ref(0);

const showNewTopic = ref(false);
const newTopic = ref({ title: "", body: "" });
const topicLoading = ref(false);
const topicError = ref("");

const openedTopic = ref(null);
const openedTopicLoading = ref(false);
const openedTopicError = ref("");
const replyBody = ref("");
const replyLoading = ref(false);
const replyingTo = ref(null);

const editTopic = ref(false);
const editTitle = ref("");
const editPost = ref(null);
const editBody = ref("");

const showReport = ref(false);
const reportTarget = ref(null);

const canCreateTopic = computed(() => auth.hasPermission("create_topic"));
const canCreatePost = computed(() => auth.hasPermission("create_post"));

const canEditTopic = computed(() => {
    if (!openedTopic.value) return false;
    return openedTopic.value.id_author === auth.user?.id || auth.isAdmin;
});
const canDeleteTopic = computed(() => {
    if (!openedTopic.value) return false;
    return openedTopic.value.id_author === auth.user?.id;
});

function canEditPost(post) {
    return post.id_author === auth.user?.id || auth.isAdmin;
}
function canDeletePost(post) {
    return post.id_author === auth.user?.id;
}
function canReportContent(authorId) {
    return authorId !== auth.user?.id;
}

const topicBody = computed(() => openedTopic.value?.posts?.[0] ?? null);
const topicReplies = computed(() => openedTopic.value?.posts?.slice(1) ?? []);

function getParentPost(parentId) {
    return openedTopic.value?.posts?.find((p) => p.id === parentId) ?? null;
}

function openReport(type, id, title) {
    reportTarget.value = { type, id, title };
    showReport.value = true;
}

let debounceTimer = null;
function debouncedFetch() {
    clearTimeout(debounceTimer);
    page.value = 1;
    debounceTimer = setTimeout(fetchTopics, 300);
}

function changePage(p) {
    page.value = p;
    window.scrollTo({ top: 0, behavior: "smooth" });
}

async function fetchTopics() {
    loading.value = true;
    try {
        const params = { page: page.value, limit: LIMIT };
        if (search.value) params.search = search.value;
        const { data } = await api.get("/forum/topics", { params });
        topics.value = data.data ?? [];
        total.value = data.total ?? 0;
    } catch (err) {
        console.error('fetchTopics error:', err)
        topics.value = [];
        total.value = 0;
    } finally {
        loading.value = false;
    }
}

async function openTopic(topic) {
    openedTopic.value = null;
    openedTopicError.value = "";
    replyBody.value = "";
    replyingTo.value = null;
    openedTopicLoading.value = true;
    try {
        const { data } = await api.get(`/forum/topics/${topic.id}`);
        openedTopic.value = data;
    } catch {
        openedTopicError.value = "Impossible de charger ce sujet.";
    } finally {
        openedTopicLoading.value = false;
    }
}

function closeTopic() {
    openedTopic.value = null;
    openedTopicLoading.value = false;
    openedTopicError.value = "";
    replyingTo.value = null;
    replyBody.value = "";
    editTopic.value = false;
    editTitle.value = "";
    editPost.value = null;
    editBody.value = "";
}

function startReply(post) {
    replyingTo.value = post;
}

function cancelReply() {
    replyingTo.value = null;
}

function startEditTopic() {
    editTitle.value = openedTopic.value?.title ?? "";
    editTopic.value = true;
}

function startEditPost(post) {
    editPost.value = post.id;
    editBody.value = post.body;
}

async function createTopic() {
    if (!newTopic.value.title.trim() || !newTopic.value.body.trim()) {
        topicError.value = "Le titre et le message sont requis.";
        return;
    }
    topicLoading.value = true;
    topicError.value = "";
    try {
        await api.post("/forum/topics", newTopic.value);
        page.value = 1;
        await fetchTopics();
        showNewTopic.value = false;
        newTopic.value = { title: "", body: "" };
    } catch (e) {
        topicError.value =
            e.response?.data?.error ?? "Une erreur est survenue.";
    } finally {
        topicLoading.value = false;
    }
}

async function saveTopicTitle() {
    if (!editTitle.value.trim() || !openedTopic.value) return;
    try {
        const { data } = await api.patch(
            `/forum/topics/${openedTopic.value.id}`,
            { title: editTitle.value.trim() }
        );
        openedTopic.value = data;
        const t = topics.value.find((t) => t.id === data.id);
        if (t) t.title = data.title;
        editTopic.value = false;
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors de la modification du titre.')
    }
}

async function deleteTopic() {
    if (!openedTopic.value) return;
    const id = openedTopic.value.id;
    try {
        await api.delete(`/forum/topics/${id}`);
        topics.value = topics.value.filter((t) => t.id !== id);
        total.value = Math.max(total.value - 1, 0);
        closeTopic();
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors de la suppression du sujet.')
    }
}

async function savePostBody(post) {
    if (!editBody.value.trim() || !openedTopic.value) return;
    try {
        const { data } = await api.patch(
            `/forum/topics/${openedTopic.value.id}/posts/${post.id}`,
            { body: editBody.value.trim() }
        );
        openedTopic.value = data;
        editPost.value = null;
        editBody.value = "";
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors de la modification du message.')
    }
}

async function deletePost(post) {
    if (!openedTopic.value) return;
    try {
        const { data } = await api.delete(
            `/forum/topics/${openedTopic.value.id}/posts/${post.id}`
        );
        openedTopic.value = data;
        const t = topics.value.find((t) => t.id === data.id);
        if (t) t.replies_count = Math.max((data.posts?.length ?? 1) - 1, 0);
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors de la suppression du message.')
    }
}

async function submitReply() {
    if (!replyBody.value.trim() || !openedTopic.value) return;
    replyLoading.value = true;
    try {
        const payload = { body: replyBody.value };
        if (replyingTo.value) payload.id_parent_post = replyingTo.value.id;
        const { data } = await api.post(
            `/forum/topics/${openedTopic.value.id}/posts`,
            payload
        );
        openedTopic.value = data;
        replyBody.value = "";
        replyingTo.value = null;
        const t = topics.value.find((t) => t.id === data.id);
        if (t) t.replies_count = Math.max((data.posts?.length ?? 1) - 1, 0);
    } catch (err) {
        alert(err?.response?.data?.error ?? 'Erreur lors de l\'envoi de la réponse.')
    } finally {
        replyLoading.value = false;
    }
}

watch(page, fetchTopics);

onMounted(fetchTopics);
</script>
