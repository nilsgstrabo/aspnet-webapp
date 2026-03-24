using Azure.Security.KeyVault.Secrets;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;

namespace aspnet_webapp.Services {
    public interface IUserInfoService
    {
        Task<UserInfo> GetUserInfo();
    }

    public class UserInfo
    {
        public string Name { get; set; }

    }

    public class UserInfoService : IUserInfoService
    {
        
        private readonly SecretClient _secretClient;

        public UserInfoService(SecretClient secretClient)
        {
            _secretClient=secretClient;
        }

        public async Task<UserInfo> GetUserInfo()
        {
            await Task.CompletedTask;
            return new UserInfo {Name="Nils"};
        }


    }

    
}