<template>
  <div class="profile-page">
    <div class="profile-container">
      <div class="profile-header">
        <div class="profile-avatar">
          <div class="avatar-icon">👤</div>
        </div>
        <h1>Мой профиль</h1>
        <p>Управление личными данными</p>
      </div>

      <div class="profile-content">
        <div class="info-card">
          <div class="card-header">
            <h3>📋 Информация об аккаунте</h3>
            <span class="member-since"
              >Участник с {{ formatDate(userInfo.created_at) }}</span
            >
          </div>

          <div class="info-item">
            <span class="info-label">📧 Email:</span>
            <span class="info-value">{{ userInfo.email }}</span>
          </div>

          <div class="info-item">
            <span class="info-label">👤 Имя пользователя:</span>
            <span class="info-value">{{ userInfo.username }}</span>
          </div>

          <div class="info-item">
            <span class="info-label">📅 Дата регистрации:</span>
            <span class="info-value">{{
              formatDate(userInfo.created_at)
            }}</span>
          </div>

          <div class="info-item">
            <span class="info-label">🆔 ID пользователя:</span>
            <span class="info-value">#{{ userInfo.id }}</span>
          </div>
        </div>

        <div class="stats-card">
          <h3>📊 Моя статистика</h3>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-value">{{ rentalsCount }}</div>
              <div class="stat-label">Всего аренд</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ activeRentalsCount }}</div>
              <div class="stat-label">Активных аренд</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ totalSpent }} ₽</div>
              <div class="stat-label">Потрачено всего</div>
            </div>
          </div>
        </div>

        <div class="edit-card">
          <h3>✏️ Редактирование профиля</h3>

          <form @submit.prevent="updateProfile">
            <div class="form-group">
              <label>👤 Имя пользователя</label>
              <input
                type="text"
                v-model="editForm.username"
                :placeholder="userInfo.username"
                required
                minlength="3"
                maxlength="50"
              />
            </div>

            <div class="form-group">
              <label>📧 Email</label>
              <input
                type="email"
                v-model="editForm.email"
                :placeholder="userInfo.email"
                required
              />
            </div>

            <div class="form-group">
              <label>🔒 Текущий пароль (обязательно для сохранения)</label>
              <input
                type="password"
                v-model="editForm.current_password"
                placeholder="Введите текущий пароль"
                required
              />
            </div>

            <div class="form-group">
              <label>🔒 Новый пароль (оставьте пустым, чтобы не менять)</label>
              <input
                type="password"
                v-model="editForm.new_password"
                placeholder="Новый пароль (мин. 6 символов)"
              />
            </div>

            <div class="form-group">
              <label>🔒 Подтверждение нового пароля</label>
              <input
                type="password"
                v-model="editForm.confirm_password"
                placeholder="Подтвердите новый пароль"
              />
            </div>

            <div class="form-actions">
              <button type="submit" class="btn-save" :disabled="loading">
                {{ loading ? "Сохранение..." : "Сохранить изменения" }}
              </button>
              <button type="button" class="btn-cancel" @click="resetForm">
                Отмена
              </button>
            </div>
          </form>
        </div>

        <div class="security-note">
          <div class="note-icon">🔒</div>
          <div class="note-text">
            <strong>Безопасность:</strong> Ваш пароль хранится в зашифрованном
            виде. Никогда не передавайте свои данные третьим лицам.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";

export default {
  name: "ProfilePage",
  setup() {
    const router = useRouter();
    const loading = ref(false);
    const loadingProfile = ref(true);
    const userInfo = ref({
      id: null,
      username: "",
      email: "",
      created_at: "",
    });
    const rentals = ref([]);
    let isComponentMounted = true;

    const editForm = ref({
      username: "",
      email: "",
      current_password: "",
      new_password: "",
      confirm_password: "",
    });

    const rentalsCount = ref(0);
    const activeRentalsCount = ref(0);
    const totalSpent = ref(0);

    // Загрузка данных пользователя
    const fetchUserInfo = async () => {
      const token = localStorage.getItem("auth_token");
      if (!token) {
        router.push("/");
        return;
      }

      try {
        const response = await axios.get("/api/user/profile");
        if (isComponentMounted) {
          userInfo.value = response.data;
          editForm.value.username = response.data.username;
          editForm.value.email = response.data.email;
        }
      } catch (error) {
        console.error("Error fetching user info:", error);
        if (error.response?.status === 401) {
          localStorage.removeItem("auth_token");
          router.push("/");
        }
      }
    };

    // Загрузка аренд пользователя
    const fetchUserRentals = async () => {
      try {
        const response = await axios.get("/api/my-rentals");
        if (isComponentMounted) {
          const userRentals = Array.isArray(response.data) ? response.data : [];
          rentals.value = userRentals;
          rentalsCount.value = userRentals.length;
          activeRentalsCount.value = userRentals.filter(
            (r) => r.status === "active",
          ).length;
          totalSpent.value = userRentals.reduce(
            (sum, rental) => sum + (rental.total_price || 0),
            0,
          );
        }
      } catch (error) {
        console.error("Error fetching rentals:", error);
        if (isComponentMounted) {
          rentals.value = [];
        }
      } finally {
        if (isComponentMounted) {
          loadingProfile.value = false;
        }
      }
    };

    // Обновление профиля
    const updateProfile = async () => {
      // Валидация
      if (!editForm.value.current_password) {
        alert("Введите текущий пароль для подтверждения изменений");
        return;
      }

      if (editForm.value.new_password) {
        if (editForm.value.new_password.length < 6) {
          alert("Новый пароль должен содержать минимум 6 символов");
          return;
        }
        if (editForm.value.new_password !== editForm.value.confirm_password) {
          alert("Новый пароль и подтверждение не совпадают");
          return;
        }
      }

      loading.value = true;
      try {
        const updateData = {
          username: editForm.value.username,
          email: editForm.value.email,
          current_password: editForm.value.current_password,
          new_password: editForm.value.new_password || undefined,
        };

        const response = await axios.put("/api/user/profile", updateData);

        if (response.data.token) {
          localStorage.setItem("auth_token", response.data.token);
        }

        alert("Профиль успешно обновлён!");

        await fetchUserInfo();
        await fetchUserRentals();

        editForm.value.current_password = "";
        editForm.value.new_password = "";
        editForm.value.confirm_password = "";
      } catch (error) {
        console.error("Error updating profile:", error);
        if (error.response?.status === 401) {
          alert("Сессия истекла. Пожалуйста, войдите заново.");
          localStorage.removeItem("auth_token");
          router.push("/");
        } else {
          alert(error.response?.data?.error || "Ошибка при обновлении профиля");
        }
      } finally {
        loading.value = false;
      }
    };

    const resetForm = () => {
      editForm.value.username = userInfo.value.username;
      editForm.value.email = userInfo.value.email;
      editForm.value.current_password = "";
      editForm.value.new_password = "";
      editForm.value.confirm_password = "";
    };

    const formatDate = (date) => {
      if (!date) return "Дата не указана";
      try {
        return new Date(date).toLocaleDateString("ru-RU", {
          day: "numeric",
          month: "long",
          year: "numeric",
        });
      } catch (e) {
        return "Дата не указана";
      }
    };

    onMounted(() => {
      isComponentMounted = true;
      fetchUserInfo();
      fetchUserRentals();
    });

    onBeforeUnmount(() => {
      isComponentMounted = false;
    });

    return {
      userInfo,
      rentals,
      editForm,
      loading,
      loadingProfile,
      rentalsCount,
      activeRentalsCount,
      totalSpent,
      updateProfile,
      resetForm,
      formatDate,
    };
  },
};
</script>
<style scoped>
.profile-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8f0fe 100%);
  padding: 2rem;
}

.profile-container {
  max-width: 900px;
  margin: 0 auto;
}

.profile-header {
  text-align: center;
  margin-bottom: 2rem;
}

.profile-avatar {
  margin-bottom: 1rem;
}

.avatar-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  color: white;
  box-shadow: 0 5px 20px rgba(0, 102, 204, 0.3);
}

.profile-header h1 {
  font-size: 2rem;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 0.5rem;
}

.profile-header p {
  color: #718096;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-card,
.stats-card,
.edit-card {
  background: white;
  border-radius: 20px;
  padding: 1.5rem;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 102, 204, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.card-header h3 {
  color: #2d3748;
  margin: 0;
}

.member-since {
  color: #0066cc;
  font-size: 0.85rem;
  background: #f0f7ff;
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
}

.info-item {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid #e2e8f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  width: 180px;
  font-weight: 500;
  color: #4a5568;
}

.info-value {
  flex: 1;
  color: #2d3748;
  font-weight: 500;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.stat-item {
  text-align: center;
  padding: 1rem;
  background: #f7fafc;
  border-radius: 15px;
  transition: all 0.3s;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.1);
}

.stat-value {
  font-size: 1.8rem;
  font-weight: bold;
  color: #0066cc;
}

.stat-label {
  color: #718096;
  font-size: 0.85rem;
  margin-top: 0.25rem;
}

.edit-card h3 {
  color: #2d3748;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4a5568;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
}

.form-group input:disabled {
  background: #f7fafc;
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-save,
.btn-cancel {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-save {
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  flex: 2;
}

.btn-save:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
}

.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-cancel {
  background: #e2e8f0;
  color: #4a5568;
  flex: 1;
}

.btn-cancel:hover {
  background: #cbd5e0;
  transform: translateY(-2px);
}

.security-note {
  background: #fff5e6;
  border-radius: 15px;
  padding: 1rem;
  display: flex;
  gap: 1rem;
  align-items: center;
  border: 1px solid #fed7aa;
}

.note-icon {
  font-size: 1.5rem;
}

.note-text {
  flex: 1;
  font-size: 0.9rem;
  color: #4a5568;
}

.note-text strong {
  color: #ed8936;
}

@media (max-width: 768px) {
  .profile-page {
    padding: 1rem;
  }

  .profile-header h1 {
    font-size: 1.5rem;
  }

  .avatar-icon {
    width: 60px;
    height: 60px;
    font-size: 2rem;
  }

  .info-item {
    flex-direction: column;
    gap: 0.25rem;
  }

  .info-label {
    width: auto;
  }

  .card-header {
    flex-direction: column;
    text-align: center;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn-save,
  .btn-cancel {
    width: 100%;
  }
}
</style>
